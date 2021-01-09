package aws

import (
	"encoding/json"
	"errors"

	"github.com/janeczku/go-spinner"
	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

const (
	operatorApiName = "operator-api-gateway"
)

func setRestApiID(cfg *config.TemplateConfig) error {
	if cfg.RestApiID != "" {
		return nil
	}

	// Look for existing REST APIs
	apis, operatorApiExists, err := getRestApis()
	if err != nil {
		return err
	}

	var restApiID string
	if len(apis) == 0 {
		// Create a new rest API
		restApiID, err = createRestApi()
		if err != nil {
			return err
		}
	} else {
		// Allow the user to create a new API gateway
		// if the operator one doesn't alredy exist
		restApiID, err := command.PromptForValue("AWS API Gateway", apis, !operatorApiExists)
		if err != nil {
			return err
		}
		if restApiID == "" {
			restApiID, err = createRestApi()
			if err != nil {
				return err
			}
		}
	}

	cfg.RestApiID = restApiID
	viper.Set(config.RestApiID, cfg.RestApiID)
	return nil
}

func setRestApiRootResourceID(cfg *config.TemplateConfig) error {
	if cfg.RestApiRootID != "" {
		return nil
	}

	s := spinner.StartNew("Collecting API root resource ID...")
	defer s.Stop()
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id", cfg.RestApiID,
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

func getRestApis() (map[string]string, bool, error) {
	s := spinner.StartNew("Collecting AWS REST APIs...")
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
		if apiGateway.Name == operatorApiName {
			operatorApiGatewayExists = true
		}
	}
	return apiGatewayIDs, operatorApiGatewayExists, nil
}

func createRestApi() (string, error) {
	s := spinner.StartNew("Creating an AWS REST API...")
	defer s.Stop()
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"create-rest-api",
		"--name", operatorApiName,
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
	if err := deployRestApi(result.ApiID); err != nil {
		return "", err
	}
	return result.ApiID, nil
}

func deployRestApi(apiID string) error {
	s := spinner.StartNew("Deploying the AWS REST API...")
	defer s.Stop()
	return command.Execute("aws", []string{
		"apigateway",
		"create-deployment",
		"--rest-api-id", apiID,
		"--stage-name", "prod",
	}, true)
}
