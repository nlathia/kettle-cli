package clouds

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
)

type GoogleCloudRun struct{}

func (g GoogleCloudRun) Deploy(directory string, cfg *config.TemplateConfig) error {
	projectID := viper.GetString(config.ProjectID)
	if projectID == "" {
		return errors.New("please run operator init")
	}

	// Build the docker image
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	commandArgs := []string{
		"builds",
		"submit",
		"--tag", fmt.Sprintf("gcr.io/%s/%s", projectID, cfg.DirectoryName),
	}
	fmt.Println("🏭  Building: ", cfg.DirectoryName, fmt.Sprintf("as a %s function", cfg.Type))
	err := executeCommand("gcloud", commandArgs)
	if err != nil {
		return err
	}

	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	commandArgs = []string{
		"run",
		"deploy",
		cfg.DirectoryName, // The cloud run service is named the same as the directory
		"--image", fmt.Sprintf("gcr.io/%s/%s", projectID, cfg.DirectoryName),
		"--platform", "managed",
		"--allow-unauthenticated",
		"--region=europe-west2",
	}
	fmt.Println("🚢  Deploying ", cfg.DirectoryName, fmt.Sprintf("as a %s function", cfg.Type))
	return executeCommand("gcloud", commandArgs)
}
