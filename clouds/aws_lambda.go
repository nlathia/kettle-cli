package clouds

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/operatorai/operator/config"
)

const (
	deploymentPackage = "deployment.zip"
)

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Setup() error {
	// @TODO: enable selecting whether to create .zip or image-based lambdas
	// @TODO: aws iam create-role
	// https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-awscli.html
	return nil
}

func (AWSLambdaFunction) Deploy(directory string, config *config.TemplateConfig) error {
	// Remove any existing deployment package
	if err := removeExistingDeployment(); err != nil {
		return err
	}

	// Store the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create the path to the zip file
	deploymentFile := path.Join(rootDir, deploymentPackage)

	// Figure out the path to the site-packages directory
	// @TODO this assumes they exist, without checking if they really do
	sitePackages, err := getPyenvSitePackagesDirectory()
	if err != nil {
		return err
	}

	// Change to the directory where the site-packages are stored
	// So that we can add them to the zip file as a directory
	os.Chdir(sitePackages)

	// Build the zip file, starting with the site-packages/
	fmt.Println(fmt.Sprintf("🧱  Building deployment archive: %s", sitePackages))
	err = executeCommand("zip", []string{
		"-r",
		deploymentFile,
		".",
	})
	if err != nil {
		return err
	}

	// Change back to the root directory
	// So that we can add its contents to the zip file
	os.Chdir(rootDir)
	fmt.Println(fmt.Sprintf("🧱  Building deployment archive: %s", rootDir))
	err = executeCommand("zip", []string{
		"-g",
		deploymentPackage,
		"-r",
		".",
	})
	if err != nil {
		return err
	}

	if lambdaExists(config.Name) {
		// Update the existing function
		fmt.Println("🚢  Updating ", config.Name, ", an existing AWS Lambda function")
		err = executeCommand("aws", []string{
			"lambda",
			"update-function-code",
			"--function-name", config.Name,
			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
		})
		if err != nil {
			return err
		}
	} else {
		// Create the function for the first time
		fmt.Println("🚢  Deploying ", config.Name, "as a new AWS Lambda function")
		fmt.Println("⏭  Entry point: ", config.FunctionName, fmt.Sprintf("(%s)", config.Runtime))
		err = executeCommand("aws", []string{
			"lambda",
			"create-function",
			"--function-name", config.Name,
			"--runtime", config.Runtime,
			"--role", "@TODO",
			"--handler", fmt.Sprintf("main.%s", config.FunctionName),
			"--package-type", "Zip",
			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
			// "--timeout", <value>,
			// "--memory-size", <value>,
		})
		if err != nil {
			return err
		}
	}

	// aws lambda wait function-active --function-name config.Name
	return nil
}

func getPyenvSitePackagesDirectory() (string, error) {
	pyenvRoot, err := executeCommandWithResult("pyenv", []string{"root"})
	if err != nil {
		return "", err
	}

	pyenvLocal, err := executeCommandWithResult("pyenv", []string{"local"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/versions/%s/lib/python3.7/site-packages/",
		strings.Trim(string(pyenvRoot), "\n"),
		strings.Trim(string(pyenvLocal), "\n"),
	), nil
}

// removeExistingDeployment removes any existing deployment zip file
// if it exits
func removeExistingDeployment() error {
	if _, err := os.Stat(deploymentPackage); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(deploymentPackage)
}

// lambdaExists queries whether a lambda function already exists
// any error is
func lambdaExists(name string) bool {
	fmt.Println("🚢  Checking if ", name, "already exists as a lambda function")
	err := executeCommand("aws", []string{
		"lambda",
		"get-function",
		"--function-name",
		name,
	})
	if err != nil {
		return false
	}
	return true
}
