package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

func SetAccountID(settings *config.Settings) error {
	if settings.AccountID != "" {
		return nil
	}

	output, err := command.ExecuteWithResult("aws", []string{
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

	settings.AccountID = result.Account
	return nil
}
