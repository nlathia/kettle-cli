package clouds

import (
	"fmt"

	"github.com/operatorai/operator/config"
)

type GoogleCloudRun struct{}

func (g GoogleCloudRun) Deploy(directory string, config *config.TemplateConfig) error {
	// gcloud config get-value project
	// @TODO we should not need to call this on every build
	projectID, err := getGoogleCloudProject()
	if err != nil {
		return err
	}

	// Build the docker image
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	commandArgs := []string{
		"builds",
		"submit",
		"--tag", fmt.Sprintf("gcr.io/%s/%s", projectID, config.DirectoryName),
	}
	fmt.Println("🏭  Building: ", config.DirectoryName, fmt.Sprintf("as a %s function", config.Type))
	err = executeCommand("gcloud", commandArgs)
	if err != nil {
		return err
	}

	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	commandArgs = []string{
		"run",
		"deploy",
		config.DirectoryName, // The cloud run service is named the same as the directory
		"--image", fmt.Sprintf("gcr.io/%s/%s", projectID, config.DirectoryName),
		"--platform", "managed",
		"--allow-unauthenticated",
		"--region=europe-west2",
	}
	fmt.Println("🚢  Deploying ", config.DirectoryName, fmt.Sprintf("as a %s function", config.Type))
	return executeCommand("gcloud", commandArgs)
}
