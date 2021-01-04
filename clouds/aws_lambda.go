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
		Label: "Available AWS IAM Roles",
		Key:   config.IAMRole,

		// Flags are currently unsupported because there's no quick way
		// to validate an ARN
		// FlagKey:           "aws-iam-role",
		// FlagDescription:   "The name of the AWS IAM role to use when deploying lambdas",
		// ValidationFunc:    validateAWSRoleExists,

		CollectValuesFunc: getAWSRoles,
	},
}

func (AWSLambdaFunction) Setup() error {
	// @TODO: enable selecting whether to create .zip or image-based lambdas
	// @TODO: aws iam create-role
	// https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-awscli.html
	return preferences.Collect(AWSConfigChoices)
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
		// @TODO aws lambda wait function-updated --function-name config.Name
	} else {
		// Create the function for the first time
		// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/create-function.html
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
		// @TODO aws lambda wait function-active --function-name config.Name
		// https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/wait/index.html#cli-aws-lambda-wait
	}

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

	// @TODO suppress output from this command
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

	if len(roles) == 0 {
		s.Stop()
		prompt := promptui.Prompt{
			Label:     "No matching AWS IAM roles. Create a new one?",
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

func createIAMRole() (map[string]string, error) {
	s := spinner.StartNew("Creating AWS IAM role for lambda.amazonaws.com...")
	defer s.Stop()

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

func validateAWSRoleExists(arn string) error {
	return nil
}
