package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

func setRestApiResourceID(cfg *config.TemplateConfig) error {
	if cfg.RestApiResourceID != "" {
		return nil
	}

	// Look for existing resource ID
	resourceID, resourceHasPOSTMethod, err := getRestApiResource(cfg)
	if err != nil {
		return err
	}
	if resourceID == "" {
		// Create a resource in the API
		output, err := command.ExecuteWithResult("aws", []string{
			"apigateway",
			"create-resource",
			"--rest-api-id", cfg.RestApiID,
			"--path-part", cfg.Name,
			"--parent-id", cfg.RestApiRootID,
		})
		if err != nil {
			return err
		}

		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(output, &result); err != nil {
			return err
		}
		cfg.RestApiResourceID = result.ID
	} else {
		// Use the existing resource ID
		cfg.RestApiResourceID = resourceID
	}
	if !resourceHasPOSTMethod {
		if err := addResourcePOSTMethod(cfg); err != nil {
			return err
		}
	}
	return nil
}

func getRestApiResource(cfg *config.TemplateConfig) (string, bool, error) {
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id", cfg.RestApiID,
	})
	if err != nil {
		return "", false, err
	}

	var results struct {
		Items []struct {
			PathPart        string `json:"pathPart"`
			ID              string `json:"id"`
			ResourceMethods struct {
				POST *struct{} `json:"POST"`
			} `json:"resourceMethods"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return "", false, err
	}

	for _, result := range results.Items {
		if result.PathPart == cfg.Name {
			return result.ID, (result.ResourceMethods.POST != nil), nil
		}
	}
	return "", false, nil
}

func addResourcePOSTMethod(cfg *config.TemplateConfig) error {
	_, resourceHasPOSTMethod, err := getRestApiResource(cfg)
	if err != nil {
		return err
	}
	if resourceHasPOSTMethod {
		return nil
	}

	apiKeySetting := "--no-api-key-required"
	requireApiKey, err := command.PromptToConfirm("Require an API key to call the URL")
	if err != nil {
		return err
	}
	if requireApiKey {
		apiKeySetting = "--api-key-required"
		// Note: if an api key is required, then there is more set up to do:

		// 1. Create a usage plan and add the usage plan to the rest api (& stage)
		// aws apigateway [get-usage-plans | create-usage-plan]
		// aws apigateway create-usage-plan --name "New Usage Plan" --description "A new usage plan" --throttle burstLimit=10,rateLimit=5 --quota limit=500,offset=0,period=MONTH --stage-keys restApiId='a1b2c3d4e5',stageName='dev'
		// apiId

		// 3. Generate an API key each user and add it to a usage plan in the console
		// aws apigateway create-api-key --name 'Dev API Key' --description 'Used for development' --enabled --stage-keys restApiId='a1b2c3d4e5',stageName='dev'
	}

	// Create the method
	err = command.Execute("aws", []string{
		"apigateway",
		"put-method",
		"--rest-api-id", cfg.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--authorization-type", "NONE",
		apiKeySetting,
	})
	if err != nil {
		return err
	}

	// Set the method response to JSON
	err = command.Execute("aws", []string{
		"apigateway",
		"put-method-response",
		"--rest-api-id", cfg.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--status-code", "200",
		"--response-models", "application/json=Empty",
	})
	if err != nil {
		return err
	}
	return nil
}
