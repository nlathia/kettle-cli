package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

type RestApiResource struct {
	Path          string
	ID            string
	HasPostMethod bool
}

func getRestApiResources(cfg *config.TemplateConfig) ([]*RestApiResource, error) {
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id", cfg.RestApiID,
	})
	if err != nil {
		return nil, err
	}

	var results struct {
		Items []struct {
			Path            string `json:"path"`
			ID              string `json:"id"`
			ResourceMethods struct {
				POST *struct{} `json:"POST"`
			} `json:"resourceMethods"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	resources := []*RestApiResource{}
	for _, result := range results.Items {
		fmt.Println(fmt.Sprintf("Found resource: %s, %s", result.Path, result.ID))
		resources = append(resources, &RestApiResource{
			Path:          result.Path,
			ID:            result.ID,
			HasPostMethod: (result.ResourceMethods.POST != nil),
		})
	}
	return resources, nil
}

func getResourceWithPath(resources []*RestApiResource, pathPart string) *RestApiResource {
	for _, resource := range resources {
		if resource.Path == strings.Join([]string{"/", pathPart}, "") {
			return resource
		}
	}
	return nil
}

func setRestApiResourceID(resources []*RestApiResource, cfg *config.TemplateConfig) error {
	if cfg.RestApiResourceID != "" {
		return nil
	}

	// Look for existing resource ID
	restApiResource := getResourceWithPath(resources, cfg.Name)
	if restApiResource != nil {
		// Use the existing resource ID
		cfg.RestApiResourceID = restApiResource.ID
	} else {
		// Not found: create a resource in the API
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
		restApiResource = &RestApiResource{
			Path:          cfg.Name,
			ID:            result.ID,
			HasPostMethod: false,
		}
	}

	// Check for POST method
	if err := addResourcePOSTMethod(restApiResource, cfg); err != nil {
		return err
	}
	return nil
}

func setRestApiRootResourceID(resources []*RestApiResource, cfg *config.TemplateConfig) error {
	if cfg.RestApiRootID != "" {
		return nil
	}
	if cfg.RestApiID == "" {
		return errors.New("rest api id not set")
	}

	resource := getResourceWithPath(resources, "")
	if resource == nil {
		return errors.New("did not find root apigateway resource")
	}

	cfg.RestApiRootID = resource.ID
	viper.Set(config.RestApiRootResource, resource.ID)
	return nil
}

func addResourcePOSTMethod(resource *RestApiResource, cfg *config.TemplateConfig) error {
	if resource.HasPostMethod {
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
