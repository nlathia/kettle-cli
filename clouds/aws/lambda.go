package aws

import (
	"errors"
	"fmt"
	"strings"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🚢  Deploying ", cfg.Name, "as an AWS Lambda function")
	fmt.Println("⏭  Entry point: ", cfg.FunctionName, fmt.Sprintf("(%s)", cfg.Runtime))

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

		addToApi, err := command.PromptToConfirm("Add Lambda function to a REST API")
		if err != nil {
			return err
		}
		if addToApi {
			if err := addLambdaToRestAPI(deploymentArchive, cfg); err != nil {
				return err
			}

			url := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
				cfg.RestApiID,
				cfg.DeploymentRegion,
				cfg.Name,
			)
			fmt.Println("🔍  API Endpoint: ", url)
		}
	}

	return waitForLambda(waitType, cfg)
}

func lambdaFunctionExists(name string) (bool, error) {
	err := command.Execute("aws", []string{
		"lambda",
		"get-function",
		"--function-name", name,
	})
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
	})
}

// https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html
func addLambdaToRestAPI(deploymentArchive string, cfg *config.TemplateConfig) error {

	// Select a deployment region
	// @TODO this leads to unnecessary repetition
	if err := setDeploymentRegion(cfg); err != nil {
		return err
	}

	// Create or set the REST API
	if err := setRestApiID(cfg); err != nil {
		return err
	}
	if err := setRestApiRootResourceID(cfg); err != nil {
		return err
	}

	// Create a resource in the API & create a POST method on the resource
	// @TODO this leads to the same aws cli call as the previous function
	if err := setRestApiResourceID(cfg); err != nil {
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
	if err := setAccountID(cfg); err != nil {
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
	case strings.HasPrefix(cfg.Runtime, "python"):
		handler = fmt.Sprintf("main.%s", cfg.FunctionName)
		runtime = cfg.Runtime
	case strings.HasPrefix(cfg.Runtime, "go"):
		handler = "main"
		runtime = "go1.x"
	default:
		return errors.New(fmt.Sprintf("unknown runtime: %s", cfg.Runtime))
	}

	// Create the function
	return command.Execute("aws", []string{
		"lambda",
		"create-function",
		"--function-name", cfg.Name,
		"--runtime", runtime,
		"--role", cfg.RoleArn,
		"--handler", handler,
		"--package-type", "Zip",
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	})
}

func waitForLambda(waitType string, cfg *config.TemplateConfig) error {
	return command.Execute("aws", []string{
		"lambda",
		"wait",
		waitType,
		"--function-name", cfg.Name,
	})
}

func addFunctionIntegration(cfg *config.TemplateConfig) error {
	// Create the integration
	err := command.Execute("aws", []string{
		"apigateway",
		"put-integration",
		"--rest-api-id", cfg.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--type", "AWS",
		"--integration-http-method", "POST",
		"--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/arn:aws:lambda:%s:%s:function:%s/invocations",
			cfg.DeploymentRegion,
			cfg.DeploymentRegion,
			cfg.AccountID,
			cfg.Name,
		),
	})
	if err != nil {
		return err
	}

	// Set the integration response to JSON
	return command.Execute("aws", []string{
		"apigateway",
		"put-integration-response",
		"--rest-api-id", cfg.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--status-code", "200",
		"--response-templates", "application/json=\"\"",
	})
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
				cfg.DeploymentRegion,
				cfg.AccountID,
				cfg.RestApiID,
				permission,
				cfg.Name,
			),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
