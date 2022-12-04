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

func (GoogleCloudRun) Deploy(directory string, cfg *config.Config, stg *settings.Settings, env string) error {
	environment, err := getEnvironment(stg, env)
	if err != nil {
		return err
	}

	// @TODO check if the current build already exists

	if strings.Contains(cfg.Config.Runtime, "go") {
		_ = cli.Execute("go", []string{
			"mod",
			"init",
		}, "Running go mod init")
	}

	fmt.Printf("🏭  Building: %s as a Cloud Run container in %s (%s)\n",
		cfg.ProjectName,
		environment.ProjectName,
		env,
	)
	containerTag := fmt.Sprintf("gcr.io/%s/%s", environment.ProjectID, cfg.ProjectName)
	// Build the docker container
	// gcloud builds submit --tag gcr.io/PROJECT-ID/helloworld
	err = cli.Execute("gcloud", []string{
		"builds",
		"submit",
		"--tag", containerTag,
		"--project", environment.ProjectID,
	}, "Building docker container")
	if err != nil {
		return err
	}

	// Deploy the docker container
	// gcloud run deploy --image gcr.io/PROJECT-ID/helloworld
	fmt.Printf("🚢  Deploying: %s as a Cloud Run container in %s (%s)\n",
		cfg.ProjectName,
		environment.ProjectName,
		env,
	)
	err = cli.Execute("gcloud", []string{
		"run",
		"deploy",
		cfg.ProjectName,
		"--image", containerTag,
		"--platform", "managed",
		"--project", environment.ProjectID,
		"--allow-unauthenticated",
		fmt.Sprintf("--region=%s", environment.DeploymentRegion),
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
		"--project", environment.ProjectName,
		"--region", environment.DeploymentRegion,
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
