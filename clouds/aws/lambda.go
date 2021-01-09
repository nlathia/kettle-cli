package aws

import (
	"fmt"

	"github.com/janeczku/go-spinner"
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
		waitType, err = updateLambda(deploymentArchive, cfg)
		if err != nil {
			return err
		}
	} else {
		waitType, err = createLambdaRestAPI(deploymentArchive, cfg)
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
		cfg.RestApiID,
		cfg.DeploymentRegion,
		cfg.Name,
	)
	fmt.Println("🔍  API Endpoint: ", url)
	return waitForLambda(waitType, cfg)
}

func lambdaFunctionExists(name string) (bool, error) {
	s := spinner.StartNew(fmt.Sprintf("Checking if lambda function \"%s\" exists...", name))
	defer s.Stop()
	err := command.Execute("aws", []string{
		"lambda",
		"get-function",
		"--function-name", name,
	}, true)
	if err != nil {
		return false, err
	}
	return true, nil
}

func updateLambda(deploymentArchive string, cfg *config.TemplateConfig) (string, error) {
	s := spinner.StartNew(fmt.Sprintf("Updating lambda function code: %s", cfg.Name))
	defer s.Stop()
	err := command.Execute("aws", []string{
		"lambda",
		"update-function-code",
		"--function-name", cfg.Name,
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, false)
	if err != nil {
		return "", err
	}
	return "function-updated", nil
}

// https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html
func createLambdaRestAPI(deploymentArchive string, cfg *config.TemplateConfig) (string, error) {
	// Get the current AWS account ID
	if err := setAccountID(cfg); err != nil {
		return "", err
	}

	// Select a deployment region
	if err := setDeploymentRegion(cfg); err != nil {
		return "", err
	}

	// Select or create the execution role
	if err := setExecutionRole(cfg); err != nil {
		return "", err
	}

	// Create or set the REST API
	if err := setRestApiID(cfg); err != nil {
		return "", err
	}
	if err := setRestApiRootResourceID(cfg); err != nil {
		return "", err
	}

	// Create a resource in the API & create a POST method on the resource
	if err := setRestApiResourceID(cfg); err != nil {
		return "", err
	}
	if err := createRestApiResourceMethod(cfg); err != nil {
		return "", err
	}

	// Set the Lambda function as the destination for the POST method
	if err := addFunctionIntegration(cfg); err != nil {
		return "", err
	}
	// Grant invoke permission to the API
	if err := addInvocationPermission(cfg); err != nil {
		return "", err
	}

	return "function-active", nil
}

func createFunction(deploymentArchive string, cfg *config.TemplateConfig) error {
	s := spinner.StartNew(fmt.Sprintf("Creating new lambda function: %s", cfg.Name))
	defer s.Stop()
	return command.Execute("aws", []string{
		"lambda",
		"create-function",
		"--function-name", cfg.Name,
		"--runtime", cfg.Runtime,
		"--role", cfg.RoleArn,
		"--handler", fmt.Sprintf("main.%s", cfg.FunctionName),
		"--package-type", "Zip",
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, false)
}

func waitForLambda(waitType string, cfg *config.TemplateConfig) error {
	s := spinner.StartNew(fmt.Sprintf("Finishing up. Waiting for: %s", waitType))
	defer s.Stop()
	return command.Execute("aws", []string{
		"lambda",
		"wait",
		waitType,
		"--function-name", cfg.Name,
	}, false)
}

func addFunctionIntegration(cfg *config.TemplateConfig) error {
	s := spinner.StartNew("Integrating the API and Lambda function...")
	defer s.Stop()

	// Create the integration
	err := command.Execute("aws", []string{
		"apigateway",
		"put-integration",
		"--rest-api-id", cfg.RestApiID,
		"--resource-id", cfg.RestApiResourceID,
		"--http-method", "POST",
		"--type", "AWS",
		"--integration-http-method", "POST",
		"--uri", fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/arn:aws:lambda:%s:%s:function:LambdaFunctionOverHttps/invocations",
			cfg.DeploymentRegion,
			cfg.DeploymentRegion,
			cfg.AccountID,
		),
	}, true)
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
	}, true)
}

func addInvocationPermission(cfg *config.TemplateConfig) error {
	s := spinner.StartNew("Adding an invocation permissions to the Lambda function...")
	defer s.Stop()
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
		}, true)
		if err != nil {
			return err
		}
	}
	return nil
}
