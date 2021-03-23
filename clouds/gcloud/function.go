package gcloud

import (
	"errors"
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type GoogleCloudFunction struct{}

// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func (GoogleCloudFunction) Deploy(directory string, cfg *config.Config, stg *settings.Settings) error {
	functionName, err := getFunctionName(cfg)
	if err != nil {
		return err
	}

	fmt.Println("🚢  Deploying ", cfg.ProjectName, "as a Google Cloud function")
	fmt.Println("⏭  Entry point: ", functionName, fmt.Sprintf("(%s)", cfg.Config.Runtime))

	fmt.Println(fmt.Sprintf("🔍  https://%s-%s.cloudfunctions.net/%s",
		stg.GoogleCloud.DeploymentRegion,
		stg.GoogleCloud.ProjectID,
		cfg.ProjectName,
	))
	return cli.Execute("gcloud", []string{
		"functions",
		"deploy",
		cfg.ProjectName,
		"--runtime", cfg.Config.Runtime,
		"--trigger-http",
		fmt.Sprintf("--entry-point=%s", functionName),
		fmt.Sprintf("--region=%s", stg.GoogleCloud.DeploymentRegion),
		"--allow-unauthenticated",
	}, "Deploying Cloud Function")
}

func getFunctionName(cfg *config.Config) (string, error) {
	for _, template := range cfg.Template {
		if template.Key == "FunctionName" {
			return template.Value, nil
		}
	}
	return "", errors.New("this template has not defined a 'FunctionName'")
}
