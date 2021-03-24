package apigateway

import (
	"encoding/json"
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type RestApiResource struct {
	Path          string
	ID            string
	HasPostMethod bool
}

func SetResourceID(resources []*RestApiResource, cfg *config.Config, stg *settings.Settings) error {
	if cfg.Config.AWS.RestApiResourceID != "" {
		return nil
	}

	// Look for existing resource ID
	restApiResource := getResourceWithPath(resources, cfg.ProjectName)
	if restApiResource == nil {
		// Not found: create a resource in the API
		output, err := cli.ExecuteWithResult("aws", []string{
			"apigateway",
			"create-resource",
			"--rest-api-id", stg.AWS.RestApiID,
			"--path-part", cfg.ProjectName,
			"--parent-id", stg.AWS.RestApiRootID,
		}, fmt.Sprintf("Creating /%s API resource", cfg.ProjectName))
		if err != nil {
			return err
		}

		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(output, &result); err != nil {
			return err
		}
		restApiResource = &RestApiResource{
			Path:          cfg.ProjectName,
			ID:            result.ID,
			HasPostMethod: false,
		}
	}

	cfg.Config.AWS.RestApiResourceID = restApiResource.ID
	// Check for POST method
	if err := addResourcePOSTMethod(restApiResource, stg.AWS.RestApiID, cfg.Config.AWS.RestApiResourceID); err != nil {
		return err
	}
	return nil
}

func addResourcePOSTMethod(resource *RestApiResource, apiID, resourceID string) error {
	if resource.HasPostMethod {
		return nil
	}

	apiKeySetting := "--no-api-key-required"
	if cli.PromptToConfirm("Require an API key to call the URL") {
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
	err := cli.Execute("aws", []string{
		"apigateway",
		"put-method",
		"--rest-api-id", apiID,
		"--resource-id", resourceID,
		"--http-method", "POST",
		"--authorization-type", "NONE",
		apiKeySetting,
	}, "Adding a POST method to the API resource")
	if err != nil {
		return err
	}

	// Set the method response to JSON
	err = cli.Execute("aws", []string{
		"apigateway",
		"put-method-response",
		"--rest-api-id", apiID,
		"--resource-id", resourceID,
		"--http-method", "POST",
		"--status-code", "200",
		"--response-models", "application/json=Empty",
	}, "Setting the resource response type to JSON")
	if err != nil {
		return err
	}
	return nil
}
