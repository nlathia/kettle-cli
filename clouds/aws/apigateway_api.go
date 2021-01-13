package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

const (
	operatorApiName = "operator-apigateway"
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
		// Allow the user to create a new REST API
		// if the operator one doesn't alredy exist
		restApiID, err = command.PromptForValue("AWS REST API", apis, !operatorApiExists)
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

func getRestApis() (map[string]string, bool, error) {
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-rest-apis",
	})
	if err != nil {
		if err.Error() == "exit status 254" {
			return map[string]string{}, false, nil
		}
		return nil, false, err
	}

	var results struct {
		Items []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, false, err
	}

	restApis := map[string]string{}
	operatorApiGatewayExists := false
	for _, restApi := range results.Items {
		restApis[restApi.Name] = restApi.ID
		if restApi.Name == operatorApiName {
			operatorApiGatewayExists = true
		}
	}
	return restApis, operatorApiGatewayExists, nil
}

func createRestApi() (string, error) {
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
	return result.ApiID, nil
}

func deployRestApi(cfg *config.TemplateConfig) error {
	return command.Execute("aws", []string{
		"apigateway",
		"create-deployment",
		"--rest-api-id", cfg.RestApiID,
		"--stage-name", "prod", // @TODO add support for different stages
	})
}
