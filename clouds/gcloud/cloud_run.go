package gcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type GoogleCloudRun struct{}

func (GoogleCloudRun) Deploy(directory string, cfg *config.Config, stg *settings.Settings) error {
	env := stg.GoogleCloud.ProdProject
	if strings.Contains(cfg.Config.Runtime, "go") {
		_ = cli.Execute("go", []string{
			"mod",
			"init",
		}, "Running go mod init")
	}

	fmt.Println("🏭  Building: ", cfg.ProjectName, "as a Cloud Run container")
	containerTag := fmt.Sprintf("gcr.io/%s/%s", env.ProjectID, cfg.ProjectName)
	// Build the docker container
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	err := cli.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", containerTag,
	}, "Building docker container")
	if err != nil {
		return err
	}

	// Deploy the docker container
	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	fmt.Println("🚢  Deploying ", cfg.ProjectName, "as a Cloud Run container")
	err = cli.Execute("gcloud", []string{
		"run",
		"deploy",
		cfg.ProjectName,
		"--image", containerTag,
		"--platform", "managed",
		"--allow-unauthenticated",
		fmt.Sprintf("--region=%s", env.DeploymentRegion),
	}, "Deploying Cloud Run container")
	if err != nil {
		return err
	}

	// Get the URL
	output, err := cli.ExecuteWithResult("gcloud", []string{
		"run",
		"services",
		"describe", cfg.ProjectName,
		"--platform", "managed",
		"--region", env.DeploymentRegion,
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
