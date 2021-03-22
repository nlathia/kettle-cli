package gcloud

import (
	"fmt"

	"github.com/operatorai/kettle-cli/command"
	"github.com/operatorai/kettle-cli/config"
)

type GoogleCloudFunction struct{}

// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func (GoogleCloudFunction) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🚢  Deploying ", cfg.Name, "as a Google Cloud function")
	fmt.Println("⏭  Entry point: ", cfg.FunctionName, fmt.Sprintf("(%s)", cfg.Settings.Runtime))
	if err := SetDeploymentRegion(cfg.Settings); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("🔍  https://%s-%s.cloudfunctions.net/%s",
		cfg.Settings.DeploymentRegion,
		cfg.Settings.ProjectID,
		cfg.Name,
	))
	return command.Execute("gcloud", []string{
		"functions",
		"deploy",
		cfg.Name,
		"--runtime", cfg.Settings.Runtime,
		"--trigger-http",
		fmt.Sprintf("--entry-point=%s", cfg.FunctionName),
		fmt.Sprintf("--region=%s", cfg.Settings.DeploymentRegion),
		"--allow-unauthenticated",
	}, "Deploying Cloud Function")
}
