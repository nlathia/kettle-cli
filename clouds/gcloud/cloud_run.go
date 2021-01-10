package gcloud

import (
	"fmt"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

type GoogleCloudRun struct{}

func (GoogleCloudRun) Deploy(directory string, cfg *config.TemplateConfig) error {
	fmt.Println("🏭  Building: ", cfg.Name, "as a Cloud Run container")
	err := setProjectID(cfg)
	if err != nil {
		return err
	}
	err = setDeploymentRegion(cfg)
	if err != nil {
		return err
	}

	// Build the docker container
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	err = command.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", fmt.Sprintf("gcr.io/%s/%s", cfg.ProjectID, cfg.Name),
	}, false)
	if err != nil {
		return err
	}

	// Deploy the docker container
	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	fmt.Println("🚢  Deploying ", cfg.Name, fmt.Sprintf("as a %s function", cfg.DeploymentType))
	return command.Execute("gcloud", []string{
		"run",
		"deploy",
		cfg.Name,
		"--image", fmt.Sprintf("gcr.io/%s/%s", cfg.ProjectID, cfg.Name),
		"--platform", "managed",
		"--allow-unauthenticated",
		fmt.Sprintf("--region=%s", cfg.DeploymentRegion),
	}, false)
}
