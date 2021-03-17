package aws

import (
	"errors"
	"fmt"
	"strings"

	"github.com/operatorai/kettle/command"
	"github.com/operatorai/kettle/config"
)

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🚢  Deploying ", cfg.Name, "as an AWS Lambda function")
	fmt.Println("⏭  Entry point: ", cfg.FunctionName, fmt.Sprintf("(%s)", cfg.Settings.Runtime))

	deploymentArchive, err := createDeploymentArchive(cfg)
	if err != nil {
		return err
	}

	var waitType string
	exists, err := lambdaFunctionExists(cfg.Name)
	if err != nil {
		return err
	}
	if exists {
		waitType = "function-updated"
		if err := updateLambda(deploymentArchive, cfg); err != nil {
			return err
		}
	} else {
		// Create the Lambda function
		waitType = "function-active"
		if err := createLambdaFunction(deploymentArchive, cfg); err != nil {
			return err
		}

		// Note: if the first deployment of a function fails after the function has
		// been created, then there is currently no way to re-deploy and create the
		// REST API. This should be changed so that a deployment asks whether to add
		// a function to an API if e.g. it hasn't already been added to one
		if command.PromptToConfirm("Add Lambda function to a REST API") {
			if err := addLambdaToRestAPI(deploymentArchive, cfg); err != nil {
				return err
			}

			url := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
				cfg.Settings.RestApiID,
				cfg.Settings.DeploymentRegion,
				cfg.Name,
			)
			fmt.Println("🔍  API Endpoint: ", url)
		}
	}

	// Clean up deployment package (ignore errors)
	_ = removeDeploymentArchive(cfg)

	return waitForLambda(waitType, cfg)
}

func lambdaFunctionExists(name string) (bool, error) {
	_, err := command.ExecuteWithResult("aws", []string{
		"lambda",
		"get-function",
		"--function-name", name,
	}, "Checking status of lambda function")
	if err != nil {
		if err.Error() == "exit status 254" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func updateLambda(deploymentArchive string, cfg *config.TemplateConfig) error {
	return command.Execute("aws", []string{
		"lambda",
		"update-function-code",
		"--function-name", cfg.Name,
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, "Updating lambda function code")
}

// https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html
func addLambdaToRestAPI(deploymentArchive string, cfg *config.TemplateConfig) error {

	// Select a deployment region (if not already set)
	if err := SetDeploymentRegion(cfg.Settings); err != nil {
		return err
	}

	// Create or set the REST API
	if err := setRestApiID(cfg.Settings); err != nil {
		return err
	}

	// Collect the available resources in the API
	resources, err := getRestApiResources(cfg)
	if err != nil {
		return err
	}

	// Set the root resource ID
	if err := setRestApiRootResourceID(resources, cfg.Settings); err != nil {
		return err
	}

	// Create a resource in the API & create a POST method on the resource
	if err := setRestApiResourceID(resources, cfg); err != nil {
		return err
	}

	// Set the Lambda function as the destination for the POST method
	if err := addFunctionIntegration(cfg); err != nil {
		return err
	}

	// Deploy the API with the new resource & integration
	if err := deployRestApi(cfg); err != nil {
		return err
	}

	// Grant invoke permission to the API
	if err := addInvocationPermission(cfg); err != nil {
		return err
	}
	return nil
}

func createLambdaFunction(deploymentArchive string, cfg *config.TemplateConfig) error {
	// Get the current AWS account ID
	if err := SetAccountID(cfg.Settings); err != nil {
		return err
	}

	// Select or create the execution role
	if err := setExecutionRole(cfg); err != nil {
		return err
	}

	// The --handler option in the create-function command changes based on the
	// programming language
	var handler string
	var runtime string
	switch {
	case strings.HasPrefix(cfg.Settings.Runtime, "python"):
		handler = fmt.Sprintf("main.%s", cfg.FunctionName)
		runtime = cfg.Settings.Runtime
	case strings.HasPrefix(cfg.Settings.Runtime, "go"):
		handler = "main"
		runtime = "go1.x"
	default:
		return errors.New(fmt.Sprintf("unknown runtime: %s", cfg.Settings.Runtime))
	}

	// Create the function
	return command.Execute("aws", []string{
		"lambda",
		"create-function",
		"--function-name", cfg.Name,
		"--runtime", runtime,
		"--role", cfg.Settings.RoleArn,
		"--handler", handler,
		"--package-type", "Zip",
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, "Creating new lambda function")
}

func waitForLambda(waitType string, cfg *config.TemplateConfig) error {
	return command.Execute("aws", []string{
		"lambda",
		"wait",
		waitType,
		"--function-name", cfg.Name,
	}, "Waiting for function to be active")
}

func addFunctionIntegration(cfg *config.TemplateConfig) error {
	// Create the integration
	err := command.Execute("aws", []string{
		"apigateway",
		"put-integration",
		"--rest-api-id", cfg.Settings.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--type", "AWS",
		"--integration-http-method", "POST",
		"--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/arn:aws:lambda:%s:%s:function:%s/invocations",
			cfg.Settings.DeploymentRegion,
			cfg.Settings.DeploymentRegion,
			cfg.Settings.AccountID,
			cfg.Name,
		),
	}, "Integrating the lambda function with the API resource")
	if err != nil {
		return err
	}

	// Set the integration response to JSON
	return command.Execute("aws", []string{
		"apigateway",
		"put-integration-response",
		"--rest-api-id", cfg.Settings.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--status-code", "200",
		"--response-templates", "application/json=\"\"",
	}, "Setting the integration response to JSON")
}

func addInvocationPermission(cfg *config.TemplateConfig) error {
	// The wildcard character (*) as the stage value indicates testing only
	permissions := map[string]string{
		"test": "*",
		"prod": "prod",
	}

	for env, permission := range permissions {
		err := command.Execute("aws", []string{
			"lambda",
			"add-permission",
			"--function-name", cfg.Name,
			"--statement-id", fmt.Sprintf("operator-apigateway-%s", env),
			"--action", "lambda:InvokeFunction",
			"--principal", "apigateway.amazonaws.com",
			"--source-arn", fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/%s/POST/%s",
				cfg.Settings.DeploymentRegion,
				cfg.Settings.AccountID,
				cfg.Settings.RestApiID,
				permission,
				cfg.Name,
			),
		}, fmt.Sprintf("Setting lambda permissions for: %s", env))
		if err != nil {
			return err
		}
	}
	return nil
}
