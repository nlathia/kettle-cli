package gcloud

import (
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
)

func build(projectID, projectName, deploymentName, env string) error {
	fmt.Printf("🏭  Building: %s as a Cloud Run container in %s (%s)\n",
		deploymentName,
		projectName,
		env,
	)

	// Check that the gcloud builds api is enabled
	// cloudbuild.googleapis.com

	containerTag := fmt.Sprintf("gcr.io/%s/%s", projectID, deploymentName)
	// Build the docker container
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	return cli.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", containerTag,
		"--project", projectName,
	}, "Building docker container")
}
