package aws

import (
	"fmt"

	"github.com/operatorai/operator/config"
)

const (
	deploymentPackage = "deployment.zip"
)

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🚢  Deploying ", cfg.Name, "as an AWS Lambda function")
	fmt.Println("⏭  Entry point: ", cfg.FunctionName, fmt.Sprintf("(%s)", cfg.Runtime))

	// Set IAM execution role
	err := setExecutionRole(cfg)
	if err != nil {
		return err
	}

	// @ TODO the rest of the deployment

	return nil
}

// @TODO
// 	// Remove any existing deployment package
// 	if err := removeExistingDeployment(); err != nil {
// 		return err
// 	}

// 	// Store the current working directory before navigating away
// 	rootDir, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}

// 	// Create the zip file, starting with the contents
// 	// of the current working directory
// 	deploymentFile := path.Join(rootDir, deploymentPackage)
// 	fmt.Println(fmt.Sprintf("🧱  Building deployment archive: %s", deploymentFile))
// 	err = command.Execute("zip", []string{
// 		"-g",
// 		deploymentPackage,
// 		"-r",
// 		".",
// 	}, true)
// 	if err != nil {
// 		return err
// 	}

// 	// Figure out the path to the site-packages directory
// 	sitePackages, err := getPyenvSitePackagesDirectory()
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := os.Stat(sitePackages); !os.IsNotExist(err) {
// 		// Change to the directory where the site-packages are stored
// 		// So that we can add them to the zip file as a directory
// 		os.Chdir(sitePackages)
// 		fmt.Println(fmt.Sprintf("🧱  Adding to deployment archive: %s", sitePackages))
// 		err = command.Execute("zip", []string{
// 			"-r",
// 			deploymentFile,
// 			".",
// 		}, true)
// 		if err != nil {
// 			return err
// 		}

// 		// Return to root directory to deploy the .zip file
// 		os.Chdir(rootDir)
// 	}

// 	// Deploy will either create or update the function
// 	// and then wait for it to be updated or active
// 	var waitCommand string

// 	fmt.Println("🚢  Deploying ", config.Name, "as an AWS Lambda function")
// 	fmt.Println("⏭  Entry point: ", config.FunctionName, fmt.Sprintf("(%s)", config.Runtime))
// 	if lambdaExists(config.Name) {
// 		// Update the existing function
// 		waitCommand = "function-updated"
// 		err = command.Execute("aws", []string{
// 			"lambda",
// 			"update-function-code",
// 			"--function-name", config.Name,
// 			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
// 		}, false)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		// Create the function for the first time
// 		// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/create-function.html
// 		waitCommand = "function-active"
// 		err = command.Execute("aws", []string{
// 			"lambda",
// 			"create-function",
// 			"--function-name", config.Name,
// 			"--runtime", config.Runtime,
// 			"--role", config.RoleArn,
// 			"--handler", fmt.Sprintf("main.%s", config.FunctionName),
// 			"--package-type", "Zip",
// 			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
// 			// "--timeout", <value>,
// 			// "--memory-size", <value>,
// 		}, false)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/wait/index.html#cli-aws-lambda-wait
// 	return command.Execute("aws", []string{
// 		"lambda",
// 		"wait",
// 		waitCommand,
// 		"--function-name",
// 		config.Name,
// 	}, false)
// }

// func getPyenvSitePackagesDirectory() (string, error) {
// 	pyenvRoot, err := command.ExecuteWithResult("pyenv", []string{"root"})
// 	if err != nil {
// 		return "", err
// 	}

// 	pyenvLocal, err := command.ExecuteWithResult("pyenv", []string{"local"})
// 	if err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("%s/versions/%s/lib/python3.7/site-packages/",
// 		strings.Trim(string(pyenvRoot), "\n"),
// 		strings.Trim(string(pyenvLocal), "\n"),
// 	), nil
// }

// // removeExistingDeployment removes the deployment.zip file, if present
// func removeExistingDeployment() error {
// 	if _, err := os.Stat(deploymentPackage); err != nil {
// 		if os.IsNotExist(err) {
// 			return nil
// 		}
// 		return err
// 	}
// 	return os.Remove(deploymentPackage)
// }

// // lambdaExists queries whether a lambda function already exists
// func lambdaExists(name string) bool {
// 	s := spinner.StartNew(fmt.Sprintf("Checking if: %s exists...", name))
// 	defer s.Stop()

// 	err := command.Execute("aws", []string{
// 		"lambda",
// 		"get-function",
// 		"--function-name",
// 		name,
// 	}, true)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }
