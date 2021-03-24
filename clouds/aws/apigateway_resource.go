package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

type RestApiResource struct {
	Path          string
	ID            string
	HasPostMethod bool
}

func getRestApiResources(stg *settings.Settings) ([]*RestApiResource, error) {
	output, err := cli.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id", stg.AWS.RestApiID,
	}, "Collecting API resources")
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

func setRestApiResourceID(resources []*RestApiResource, stg *settings.Settings) error {
	if stg.AWS.RestApiResourceID != "" {
		return nil
	}

	// Look for existing resource ID
	restApiResource := getResourceWithPath(resources, cfg.Name)
	if restApiResource != nil {
		// Use the existing resource ID
		cfg.RestApiResourceID = restApiResource.ID
	} else {
		// Not found: create a resource in the API
		output, err := cli.ExecuteWithResult("aws", []string{
			"apigateway",
			"create-resource",
			"--rest-api-id", cfg.Settings.RestApiID,
			"--path-part", cfg.Name,
			"--parent-id", cfg.Settings.RestApiRootID,
		}, fmt.Sprintf("Creating /%s API resource", cfg.Name))
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

func setRestApiRootResourceID(resources []*RestApiResource, stg *settings.Settings) error {
	if stg.AWS.RestApiRootID != "" {
		return nil
	}
	if stg.AWS.RestApiID == "" {
		return errors.New("rest api id not set")
	}

	resource := getResourceWithPath(resources, "")
	if resource == nil {
		return errors.New("did not find root apigateway resource")
	}

	stg.AWS.RestApiRootID = resource.ID
	return nil
}

func addResourcePOSTMethod(resource *RestApiResource, stg *settings.Settings) error {
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
		"--rest-api-id", stg.AWS.RestApiID,
		"--resource-id", stg.AWS.RestApiResourceID, // @TODO
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
		"--rest-api-id", stg.AWS.RestApiID,
		"--resource-id", stg.AWS.RestApiResourceID, // @TODO
		"--http-method", "POST",
		"--status-code", "200",
		"--response-models", "application/json=Empty",
	}, "Setting the resource response type to JSON")
	if err != nil {
		return err
	}
	return nil
}
