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
	return command.Execute("gcloud", []string{
		"functions",
		"deploy",
		cfg.Name, // The cloud function is named the same as the directory
		"--runtime", cfg.Runtime,
		"--trigger-http", // We only currently support http triggers
		fmt.Sprintf("--entry-point=%s", cfg.FunctionName),
		fmt.Sprintf("--region=%s", cfg.DeploymentRegion),
		"--allow-unauthenticated",
	}, false)
}
