package gcloud

import (
	"fmt"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

type GoogleCloudFunction struct{}

// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func (GoogleCloudFunction) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🚢  Deploying ", cfg.Name, "as a Google Cloud function")
	fmt.Println("⏭  Entry point: ", cfg.FunctionName, fmt.Sprintf("(%s)", cfg.Runtime))
	setDeploymentRegion(cfg)

	return command.Execute("gcloud", []string{
		"functions",
		"deploy",
		cfg.Name,
		"--runtime", cfg.Runtime,
		"--trigger-http",
		fmt.Sprintf("--entry-point=%s", cfg.FunctionName),
		fmt.Sprintf("--region=%s", cfg.DeploymentRegion),
		"--allow-unauthenticated",
	}, false)
}
