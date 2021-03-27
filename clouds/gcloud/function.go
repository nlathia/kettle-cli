package gcloud

import (
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type GoogleCloudFunction struct{}

// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func (GoogleCloudFunction) Deploy(directory string, cfg *config.Config, stg *settings.Settings) error {
	fmt.Println("🚢  Deploying ", cfg.ProjectName, "as a Google Cloud function")
	fmt.Println("⏭  Entry point: ", cfg.Config.EntryFunction, fmt.Sprintf("(%s)", cfg.Config.Runtime))

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
		fmt.Sprintf("--entry-point=%s", cfg.Config.EntryFunction),
		fmt.Sprintf("--region=%s", stg.GoogleCloud.DeploymentRegion),
		"--allow-unauthenticated",
	}, "Deploying Cloud Function")
}
