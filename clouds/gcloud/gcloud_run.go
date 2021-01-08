package gcloud

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

type GoogleCloudRun struct{}

func (GoogleCloudRun) GetConfig() *config.TemplateConfig {
	return nil
}

func (GoogleCloudRun) Setup() error {
	// return preferences.Collect(GCPConfigChoices)
	return nil
}

func (GoogleCloudRun) Deploy(directory string, cfg *config.TemplateConfig) error {
	projectID := viper.GetString(config.ProjectID)
	if projectID == "" {
		return errors.New("please run operator init")
	}

	// Build the docker image
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	fmt.Println("🏭  Building: ", cfg.Name, "as a Cloud Run container")
	err := command.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", fmt.Sprintf("gcr.io/%s/%s", projectID, cfg.Name),
	}, false)
	if err != nil {
		return err
	}

	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	fmt.Println("🚢  Deploying ", cfg.Name, fmt.Sprintf("as a %s function", cfg.DeploymentType))
	return command.Execute("gcloud", []string{
		"run",
		"deploy",
		cfg.Name, // The cloud run service is named the same as the directory
		"--image", fmt.Sprintf("gcr.io/%s/%s", projectID, cfg.Name),
		"--platform", "managed",
		"--allow-unauthenticated",
		"--region=europe-west2",
	}, false)
}
