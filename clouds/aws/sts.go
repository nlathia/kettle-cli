package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

func setAccountID(cfg *config.TemplateConfig) error {
	if cfg.AccountID != "" {
		return nil
	}

	// fmt.Println("Retrieving AWS identity...")
	// s := spinner.StartNew("Querying...")
	// defer s.Stop()
	output, err := command.ExecuteWithResult("aws", []string{
		"sts",
		"get-caller-identity",
		"--output", "json",
	})
	if err != nil {
		return err
	}

	var result struct {
		Account string `json:"Account"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}

	cfg.AccountID = result.Account
	viper.Set(config.AccountID, result.Account)
	return nil
}
