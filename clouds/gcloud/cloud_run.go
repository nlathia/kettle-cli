package gcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operatorai/kettle/command"
	"github.com/operatorai/kettle/config"
)

type GoogleCloudRun struct{}

func (GoogleCloudRun) Deploy(directory string, cfg *config.TemplateConfig) error {
	if strings.Contains(cfg.Settings.Runtime, "go") {
		_ = command.Execute("go", []string{
			"mod",
			"init",
		}, "Running go mod init")
	}

	fmt.Println("🏭  Building: ", cfg.Name, "as a Cloud Run container")
	if err := SetProjectID(cfg.Settings); err != nil {
		return err
	}
	if err := SetDeploymentRegion(cfg.Settings); err != nil {
		return err
	}

	containerTag := fmt.Sprintf("gcr.io/%s/%s", cfg.Settings.ProjectID, cfg.Name)
	// Build the docker container
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	err := command.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", containerTag,
	}, "Building docker container")
	if err != nil {
		return err
	}

	// Deploy the docker container
	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	fmt.Println("🚢  Deploying ", cfg.Name, "as a Cloud Run container")
	err = command.Execute("gcloud", []string{
		"run",
		"deploy",
		cfg.Name,
		"--image", containerTag,
		"--platform", "managed",
		"--allow-unauthenticated",
		fmt.Sprintf("--region=%s", cfg.Settings.DeploymentRegion),
	}, "Deploying Cloud Run container")
	if err != nil {
		return err
	}

	// Get the URL
	output, err := command.ExecuteWithResult("gcloud", []string{
		"run",
		"services",
		"describe", cfg.Name,
		"--platform", "managed",
		"--region", cfg.Settings.DeploymentRegion,
		"--format", "json",
	}, "Querying for Cloud Run URL")
	if err != nil {
		fmt.Println("😥  Could not retrieve URL (but the Cloud Run function has deployed)")
		return nil
	}

	var results struct {
		Status struct {
			URL string `json:"url"`
		} `json:"status"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		fmt.Println("😥  Could not parse response (but the Cloud Run function has deployed)")
		return nil
	}

	fmt.Println("🔍  API Endpoint: ", results.Status.URL)
	return nil
}
