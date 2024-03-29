package aws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

const (
	operatorExecutionRole = "operator-lambda-role"
)

func setExecutionRole(stg *settings.Settings, overwrite bool) error {
	if !overwrite {
		if stg.AWS.RoleArn != "" {
			return nil
		}
	}

	roles, operatorExecutionRoleExists, err := getExecutionRoles()
	if err != nil {
		return err
	}

	var role string
	if len(roles) == 0 {
		role, err = createExecutionRole()
		if err != nil {
			return err
		}
	} else {
		role, err = cli.PromptForValue("IAM Role", roles, !operatorExecutionRoleExists)
		if err != nil {
			return err
		}
		if role == "" {
			role, err = createExecutionRole()
			if err != nil {
				return err
			}
		}
	}

	stg.AWS.RoleArn = role
	return nil
}

func getExecutionRoles() (map[string]string, bool, error) {
	output, err := cli.ExecuteWithResult("aws", []string{
		"iam",
		"list-roles",
		"--output", "json",
	}, "Collecting available IAM roles")
	if err != nil {
		return nil, false, err
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
		return nil, false, err
	}

	operatorExecutionRoleExists := false
	roles := map[string]string{}
	for _, role := range results.Roles {
		if role.RolePolicy.Statement[0].Principal.Service == "lambda.amazonaws.com" {
			displayName := fmt.Sprintf("%s (%s)", role.RoleName, role.Path)
			roles[displayName] = role.Arn
			if role.RoleName == operatorExecutionRole {
				operatorExecutionRoleExists = true
			}
		}
	}
	return roles, operatorExecutionRoleExists, nil
}

func createExecutionRole() (string, error) {
	// Write the trust policy to a temp file
	f, err := ioutil.TempFile(".", "trust_policy*.json")
	if err != nil {
		return "", err
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
		return "", err
	}

	output, err := cli.ExecuteWithResult("aws", []string{
		"iam",
		"create-role",
		"--role-name", operatorExecutionRole,
		"--assume-role-policy-document", fmt.Sprintf("file://%s", f.Name()),
		"--output", "json",
	}, fmt.Sprintf("Creating an IAM role called: %s", operatorExecutionRole))
	if err != nil {
		return "", err
	}

	var result struct {
		Role struct {
			Arn string `json:"Arn"`
		} `json:"Role"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", err
	}
	return result.Role.Arn, nil
}
