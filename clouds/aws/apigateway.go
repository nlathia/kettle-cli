package aws

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/janeczku/go-spinner"
	"github.com/manifoldco/promptui"
	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

const (
	operatorApiGateway = "operator-api-gateway"
)

func setApiGateway(cfg *config.TemplateConfig) error {
	if cfg.RestApiID != "" {
		return nil
	}

	apis, operatorApiGatewayExists, err := getApiGateways()
	if err != nil {
		return err
	}

	var restApiID string
	if len(apis) == 0 {
		prompt := promptui.Prompt{
			Label:     "No API gateways. Create a new one",
			IsConfirm: true,
		}

		confirmed, err := prompt.Run()
		if err != nil {
			return err
		}
		if strings.ToLower(confirmed) != "y" {
			return errors.New("cancelled")
		}

		restApiID, err = createApiGateway()
		if err != nil {
			return err
		}
	} else {
		// Allow the user to create a new API gateway if the operator one
		// doesn't alredy exist
		restApiID, err := command.PromptForValue("AWS API Gateway", apis, !operatorApiGatewayExists)
		if err != nil {
			return err
		}
		if restApiID == "" {
			restApiID, err = createApiGateway()
			if err != nil {
				return err
			}
		}
	}

	cfg.RestApiID = restApiID
	viper.Set(config.RestApiID, cfg.RestApiID)
	return nil
}

func setApiGatewayRoot(cfg *config.TemplateConfig) error {
	if cfg.RestApiRootID != "" {
		return nil
	}
	if err := setApiGateway(cfg); err != nil {
		return err
	}

	s := spinner.StartNew("Collecting API root resource ID...")
	defer s.Stop()
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id",
		cfg.RestApiID,
	})
	if err != nil {
		return err
	}

	var results struct {
		Items []struct {
			Path string `json:"path"`
			ID   string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return err
	}
	if len(results.Items) == 0 {
		return errors.New("no matching apigateway resource")
	}

	cfg.RestApiRootID = results.Items[0].ID
	viper.Set(config.RestApiRootResource, results.Items[0].ID)
	return nil
}

func setApiGatewayResource(cfg *config.TemplateConfig) error {
	if cfg.RestApiResourcePath != "" {
		return nil
	}
	if err := setApiGateway(cfg); err != nil {
		return err
	}
	if err := setApiGatewayRoot(cfg); err != nil {
		return err
	}

	s := spinner.StartNew("Creating an AWS API gateway resource...")
	defer s.Stop()

	// Create a resource in the API
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"create-resource",
		"--rest-api-id",
		cfg.RestApiID,
		"--path-part",
		cfg.Name,
		"--parent-id",
		cfg.RestApiRootID,
	})
	if err != nil {
		return err
	}

	var result struct {
		Path string `json:"path"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}
	cfg.RestApiResourcePath = result.Path

	// Create POST method on the resource
	err = command.Execute("aws", []string{
		"apigateway",
		"put-method",
		"--rest-api-id",
		cfg.RestApiID,
		"--resource-id",
		result.ID,
		"--http-method",
		"POST",
		"--authorization-type",
		"NONE",
	}, true)
	if err != nil {
		return err
	}
	return nil
}

func getApiGateways() (map[string]string, bool, error) {
	s := spinner.StartNew("Collecting AWS API gateways...")
	defer s.Stop()

	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-rest-apis",
	})
	if err != nil {
		return nil, false, err
	}

	var results struct {
		Items []struct {
			ApiID       string `json:"id"`
			Name        string `json:"name"`
			CreatedDate int    `json:"createdDate"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, false, err
	}

	apiGatewayIDs := map[string]string{}
	operatorApiGatewayExists := false
	for _, apiGateway := range results.Items {
		apiGatewayIDs[apiGateway.Name] = apiGateway.ApiID
		if apiGateway.Name == operatorApiGateway {
			operatorApiGatewayExists = true
		}
	}
	return apiGatewayIDs, operatorApiGatewayExists, nil
}

func createApiGateway() (string, error) {
	s := spinner.StartNew("Creating new AWS API gateway...")
	defer s.Stop()

	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"create-rest-api",
		"--name",
	})
	if err != nil {
		return "", err
	}

	var result struct {
		ApiID string `json:"id"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}
	return result.ApiID, nil
}
