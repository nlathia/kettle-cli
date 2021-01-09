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
	if lambdaExists(cfg.Name) {
		waitType, err = updateLambda(deploymentArchive, cfg)
		if err != nil {
			return err
		}
	} else {
		waitType, err = createLambda(deploymentArchive, cfg)
		if err != nil {
			return err
		}
	}
	return waitForLambda(waitType, cfg)
}

// lambdaExists queries whether a lambda function already exists
func lambdaExists(name string) bool {
	s := spinner.StartNew(fmt.Sprintf("Checking if: %s exists...", name))
	defer s.Stop()

	err := command.Execute("aws", []string{
		"lambda",
		"get-function",
		"--function-name",
		name,
	}, true)
	if err != nil {
		return false
	}
	return true
}

func updateLambda(deploymentArchive string, cfg *config.TemplateConfig) (string, error) {
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

func createLambda(deploymentArchive string, cfg *config.TemplateConfig) (string, error) {
	err := setExecutionRole(cfg)
	if err != nil {
		return "", err
	}

	err = command.Execute("aws", []string{
		"lambda",
		"create-function",
		"--function-name", cfg.Name,
		"--runtime", cfg.Runtime,
		"--role", cfg.RoleArn,
		"--handler", fmt.Sprintf("main.%s", cfg.FunctionName),
		"--package-type", "Zip",
		"--zip-file", fmt.Sprintf("fileb://%s", deploymentArchive),
	}, false)
	if err != nil {
		return "", err
	}

	// @TODO add api gateway
	// @TODO get api gateway root ID
	return "function-active", nil
}

func waitForLambda(waitType string, cfg *config.TemplateConfig) error {
	s := spinner.StartNew(fmt.Sprintf("Deploying. Waiting for: %s", waitType))
	defer s.Stop()
	return command.Execute("aws", []string{
		"lambda",
		"wait",
		waitType,
		"--function-name",
		cfg.Name,
	}, false)
}
