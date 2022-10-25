package gcloud

import (
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type GoogleCloudFunction struct{}

// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func (GoogleCloudFunction) Deploy(directory string, cfg *config.Config, stg *settings.Settings, env string) error {
	environment, err := getEnvironment(stg, env)
	if err != nil {
		return err
	}

	fmt.Printf("🚢  Deploying %s as a Google Cloud function to %s (%s)\n",
		cfg.ProjectName,
		environment.ProjectName,
		env,
	)
	fmt.Printf("⏭  Entry point: %s (%s)\n", cfg.Config.EntryFunction, cfg.Config.Runtime)
	fmt.Printf("🔍  https://%s-%s.cloudfunctions.net/%s\n",
		environment.DeploymentRegion,
		environment.ProjectID,
		cfg.ProjectName,
	)

	return cli.Execute("gcloud", []string{
		"functions",
		"deploy",
		cfg.ProjectName,
		"--runtime", cfg.Config.Runtime,
		"--trigger-http",
		fmt.Sprintf("--entry-point=%s", cfg.Config.EntryFunction),
		fmt.Sprintf("--region=%s", environment.DeploymentRegion),
		"--allow-unauthenticated",
	}, "Deploying Cloud Function")
}
