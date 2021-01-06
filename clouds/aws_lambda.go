package clouds

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/janeczku/go-spinner"
	"github.com/manifoldco/promptui"
	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/preferences"
)

const (
	deploymentPackage = "deployment.zip"
)

type AWSLambdaFunction struct{}

var AWSConfigChoices = []*preferences.ConfigChoice{
	{
		// Pick or create an AWS IAM role for deploying Lambdas
		Label:             "Available AWS IAM Roles",
		Key:               config.IAMRole,
		FlagKey:           "aws-iam-arn",
		FlagDescription:   "The ARN of the AWS IAM role to use when deploying lambdas",
		ValidationFunc:    validateAWSRoleExists,
		CollectValuesFunc: collectAWSRoles,
	},
}

func (AWSLambdaFunction) Setup() error {
	// @TODO (Future): enable selecting whether to create .zip or image-based lambdas
	return preferences.Collect(AWSConfigChoices)
}

func (AWSLambdaFunction) Deploy(directory string, config *config.TemplateConfig) error {
	// Remove any existing deployment package
	if err := removeExistingDeployment(); err != nil {
		return err
	}

	// Store the current working directory before navigating away
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create the zip file, starting with the contents
	// of the current working directory
	deploymentFile := path.Join(rootDir, deploymentPackage)
	fmt.Println(fmt.Sprintf("🧱  Building deployment archive: %s", deploymentFile))
	err = executeCommand("zip", []string{
		"-g",
		deploymentPackage,
		"-r",
		".",
	}, true)
	if err != nil {
		return err
	}

	// Figure out the path to the site-packages directory
	sitePackages, err := getPyenvSitePackagesDirectory()
	if err != nil {
		return err
	}

	if _, err := os.Stat(sitePackages); !os.IsNotExist(err) {
		// Change to the directory where the site-packages are stored
		// So that we can add them to the zip file as a directory
		os.Chdir(sitePackages)
		fmt.Println(fmt.Sprintf("🧱  Adding to deployment archive: %s", sitePackages))
		err = executeCommand("zip", []string{
			"-r",
			deploymentFile,
			".",
		}, true)
		if err != nil {
			return err
		}

		// Return to root directory to deploy the .zip file
		os.Chdir(rootDir)
	}

	// Deploy will either create or update the function
	// and then wait for it to be updated or active
	var waitCommand string

	fmt.Println("🚢  Deploying ", config.Name, "as an AWS Lambda function")
	fmt.Println("⏭  Entry point: ", config.FunctionName, fmt.Sprintf("(%s)", config.Runtime))
	if lambdaExists(config.Name) {
		// Update the existing function
		waitCommand = "function-updated"
		err = executeCommand("aws", []string{
			"lambda",
			"update-function-code",
			"--function-name", config.Name,
			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
		}, false)
		if err != nil {
			return err
		}
	} else {
		// Create the function for the first time
		// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/create-function.html
		waitCommand = "function-active"
		err = executeCommand("aws", []string{
			"lambda",
			"create-function",
			"--function-name", config.Name,
			"--runtime", config.Runtime,
			"--role", config.IAMRole,
			"--handler", fmt.Sprintf("main.%s", config.FunctionName),
			"--package-type", "Zip",
			"--zip-file", fmt.Sprintf("fileb://%s", deploymentPackage),
			// "--timeout", <value>,
			// "--memory-size", <value>,
		}, false)
		if err != nil {
			return err
		}
	}

	// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/wait/index.html#cli-aws-lambda-wait
	return executeCommand("aws", []string{
		"lambda",
		"wait",
		waitCommand,
		"--function-name",
		config.Name,
	}, false)
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

// removeExistingDeployment removes the deployment.zip file, if present
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
func lambdaExists(name string) bool {
	s := spinner.StartNew(fmt.Sprintf("Checking if: %s exists...", name))
	defer s.Stop()

	err := executeCommand("aws", []string{
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

func getAWSRoles() (map[string]string, error) {
	s := spinner.StartNew("Collecting AWS IAM roles...")
	defer s.Stop()

	//  aws iam list-roles --output json
	output, err := executeCommandWithResult("aws", []string{
		"iam",
		"list-roles",
		"--output",
		"json",
	})
	if err != nil {
		return nil, err
	}

	var results struct {
		Roles []struct {
			RoleName   string `json:"RoleName"`
			Path       string `json:"Path"`
			Arn        string `json:"Arn"`
			RolePolicy struct {
				Statement []struct {
					Principal struct {
						Service string `json:"Service"`
					} `json:"Principal"`
				} `json:"Statement"`
			} `json:"AssumeRolePolicyDocument"`
		} `json:"Roles"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	roles := map[string]string{}
	for _, role := range results.Roles {
		if role.RolePolicy.Statement[0].Principal.Service == "lambda.amazonaws.com" {
			displayName := fmt.Sprintf("%s (%s)", role.RoleName, role.Path)
			roles[displayName] = role.Arn
		}
	}
	return roles, nil
}

func collectAWSRoles() (map[string]string, error) {
	roles, err := getAWSRoles()
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		prompt := promptui.Prompt{
			Label:     "No matching AWS IAM roles. Create a new one",
			IsConfirm: true,
		}

		confirmed, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if strings.ToLower(confirmed) == "y" {
			return createIAMRole()
		}
		return roles, errors.New("unknown input")
	}
	return roles, nil
}

func validateAWSRoleExists(arn string) error {
	roles, err := getAWSRoles()
	if err != nil {
		return err
	}

	for _, roleArn := range roles {
		if roleArn == arn {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("No matching role for ARN: %s", arn))
}

func createIAMRole() (map[string]string, error) {
	s := spinner.StartNew("Creating AWS IAM role for lambda.amazonaws.com...")
	defer s.Stop()

	// Write the trust policy to a temp file
	f, err := ioutil.TempFile(".", "trust_policy*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	trustPolicy := []byte(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`)
	if _, err = f.Write(trustPolicy); err != nil {
		return nil, err
	}

	// $ aws iam create-role --role-name lambda-ex --assume-role-policy-document file://trust-policy.json
	output, err := executeCommandWithResult("aws", []string{
		"iam",
		"create-role",
		"--role-name",
		"operator-lambda-role",
		"--assume-role-policy-document",
		fmt.Sprintf("file://%s", f.Name()),
		"--output",
		"json",
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Role struct {
			RoleName string `json:"RoleName"`
			Path     string `json:"Path"`
			Arn      string `json:"Arn"`
		} `json:"Role"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	displayName := fmt.Sprintf("%s (%s)", result.Role.RoleName, result.Role.Path)
	return map[string]string{
		displayName: result.Role.Arn,
	}, nil
}
