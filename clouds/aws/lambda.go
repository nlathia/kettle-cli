package aws

import (
	"errors"
	"fmt"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/clouds/aws/apigateway"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Deploy(directory string, cfg *config.Config, stg *settings.Settings) error {
	fmt.Println("🚢  Deploying ", cfg.ProjectName, "as an AWS Lambda function")
	fmt.Println("⏭  Entry point: ", cfg.Config.EntryFunction, fmt.Sprintf("(%s)", cfg.Config.Runtime))
	// @TODO future - container-based deployments
	deploymentArchive, err := createDeploymentArchive(cfg)
	if err != nil {
		return err
	}
	defer func() {
		// Clean up deployment package (ignore errors)
		err := removeDeploymentArchive(cfg)
		if err != nil {
			if settings.DebugMode {
				fmt.Println(err.Error())
			}
		}
	}()

	var waitType string
	exists, err := lambdaFunctionExists(cfg.ProjectName)
	if err != nil {
		return err
	}
	if exists {
		// Update the function with the new code
		waitType = "function-updated"
		if err := updateLambda(deploymentArchive, cfg); err != nil {
			return err
		}
	} else {
		// Create the Lambda function
		waitType = "function-active"
		if err := createLambdaFunction(deploymentArchive, cfg.Config.EntryFunction, cfg, stg); err != nil {
			return err
		}

		// Note: if the first deployment of a function fails after the function has
		// been created, then there is currently no way to re-deploy and create the
		// REST API. This should be changed so that a deployment asks whether to add
		// a function to an API if e.g. it hasn't already been added to one
		if cli.PromptToConfirm("Add Lambda function to a REST API") {
			if err := addLambdaToRestAPI(deploymentArchive, cfg, stg); err != nil {
				return err
			}

			url := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
				stg.AWS.RestApiID,
				stg.AWS.DeploymentRegion,
				cfg.ProjectName,
			)
			fmt.Println("🔍  API Endpoint: ", url)
		}
	}
	return waitForLambda(waitType, cfg)
}

func lambdaFunctionExists(name string) (bool, error) {
	_, err := cli.ExecuteWithResult("aws", []string{
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

func updateLambda(deploymentArchive string, cfg *config.Config) error {
	return cli.Execute("aws", []string{
		"lambda",
		"update-function-code",
		"--function-name", cfg.ProjectName,
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, "Updating lambda function code")
}

// https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html
func addLambdaToRestAPI(deploymentArchive string, cfg *config.Config, stg *settings.Settings) error {
	// Create or set the REST API
	if err := apigateway.SetRestApiID(stg); err != nil {
		return err
	}

	// Collect the available resources in the API
	resources, err := apigateway.GetResources(stg)
	if err != nil {
		return err
	}

	// Set the root resource ID
	if err := apigateway.SetRootResourceID(resources, stg); err != nil {
		return err
	}

	// Create a resource in the API & create a POST method on the resource
	if err := apigateway.SetResourceID(resources, cfg, stg); err != nil {
		return err
	}

	// Set the Lambda function as the destination for the POST method
	if err := addFunctionIntegration(cfg, stg); err != nil {
		return err
	}

	// Deploy the API with the new resource & integration
	if err := apigateway.Deploy(stg); err != nil {
		return err
	}

	// Grant invoke permission to the API
	if err := addInvocationPermission(cfg, stg); err != nil {
		return err
	}
	return nil
}

func createLambdaFunction(deploymentArchive string, functionName string, cfg *config.Config, stg *settings.Settings) error {
	// Get the current AWS account ID
	if err := SetAccountID(stg.AWS); err != nil {
		return err
	}

	// Select or create the execution role
	if err := setExecutionRole(stg); err != nil {
		return err
	}

	// The --handler option in the create-function command changes based on the
	// programming language
	var handler string
	var runtime string
	switch {
	case strings.HasPrefix(cfg.Config.Runtime, "python"):
		handler = fmt.Sprintf("main.%s", functionName)
		runtime = cfg.Config.Runtime
	case strings.HasPrefix(cfg.Config.Runtime, "go"):
		handler = "main"
		runtime = "go1.x"
	default:
		return errors.New(fmt.Sprintf("unknown runtime: %s", cfg.Config.Runtime))
	}

	// Create the function
	return cli.Execute("aws", []string{
		"lambda",
		"create-function",
		"--function-name", cfg.ProjectName,
		"--runtime", runtime,
		"--role", stg.AWS.RoleArn,
		"--handler", handler,
		"--package-type", "Zip",
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, "Creating new lambda function")
}

func waitForLambda(waitType string, cfg *config.Config) error {
	return cli.Execute("aws", []string{
		"lambda",
		"wait",
		waitType,
		"--function-name", cfg.ProjectName,
	}, "Waiting for function to be active")
}

func addFunctionIntegration(cfg *config.Config, stg *settings.Settings) error {
	// Create the integration
	err := cli.Execute("aws", []string{
		"apigateway",
		"put-integration",
		"--rest-api-id", stg.AWS.RestApiID,
		"--resource-id", cfg.Config.AWS.RestApiResourceID,
		"--http-method", "POST",
		"--type", "AWS",
		"--integration-http-method", "POST",
		"--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/arn:aws:lambda:%s:%s:function:%s/invocations",
			stg.AWS.DeploymentRegion,
			stg.AWS.DeploymentRegion,
			stg.AWS.AccountID,
			cfg.ProjectName,
		),
	}, "Integrating the lambda function with the API resource")
	if err != nil {
		return err
	}

	// Set the integration response to JSON
	return cli.Execute("aws", []string{
		"apigateway",
		"put-integration-response",
		"--rest-api-id", stg.AWS.RestApiID,
		"--resource-id", cfg.Config.AWS.RestApiResourceID,
		"--http-method", "POST",
		"--status-code", "200",
		"--response-templates", "application/json=\"\"",
	}, "Setting the integration response to JSON")
}

func addInvocationPermission(cfg *config.Config, stg *settings.Settings) error {
	// The wildcard character (*) as the stage value indicates testing only
	permissions := map[string]string{
		"test": "*",
		"prod": "prod",
	}
	for env, permission := range permissions {
		err := cli.Execute("aws", []string{
			"lambda",
			"add-permission",
			"--function-name", cfg.ProjectName,
			"--statement-id", fmt.Sprintf("operator-apigateway-%s", env),
			"--action", "lambda:InvokeFunction",
			"--principal", "apigateway.amazonaws.com",
			"--source-arn", fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/%s/POST/%s",
				stg.AWS.DeploymentRegion,
				stg.AWS.AccountID,
				stg.AWS.RestApiID,
				permission,
				cfg.ProjectName,
			),
		}, fmt.Sprintf("Setting lambda permissions for: %s", env))
		if err != nil {
			return err
		}
	}
	return nil
}
