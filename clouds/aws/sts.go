package aws

import (
	"encoding/json"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

func SetAccountID(stg *settings.AWSSettings, overwrite bool) error {
	if !overwrite {
		if stg.AccountID != "" {
			return nil
		}
	}

	output, err := cli.ExecuteWithResult("aws", []string{
		"sts",
		"get-caller-identity",
		"--output", "json",
	}, "Retrieving aws caller identity")
	if err != nil {
		return err
	}

	var result struct {
		Account string `json:"Account"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}

	stg.AccountID = result.Account
	return nil
}
