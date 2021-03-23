package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

func SetProjectID(settings *settings.GoogleCloudSettings) error {
	if settings.ProjectName != "" && settings.ProjectID != "" {
		return nil
	}

	projects, err := getGoogleCloudProjects()
	if err != nil {
		return err
	}

	projectName, projectID, err := cli.PromptForKeyValue("Google Cloud Project", projects)
	if err != nil {
		return err
	}

	settings.ProjectName = projectName
	settings.ProjectID = projectID
	return nil
}

func getGoogleCloudProjects() (map[string]string, error) {
	// gcloud projects list --format="json"
	output, err := cli.ExecuteWithResult("gcloud", []string{
		"projects",
		"list",
		"--format=\"json\"",
		fmt.Sprintf("--limit=%d", 50),
	}, "Collecting gcloud projects")
	if err != nil {
		return nil, err
	}

	var results []struct {
		ProjectID string `json:"projectId"`
		Name      string `json:"name"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	// @TODO re-enable this
	// if len(results) >= projectListLimit {
	// 	// Bail if the user has too many projects to 'reasonably' display
	// 	return nil, errors.New(fmt.Sprintf("you have %d or more Google Cloud projects. "+
	// 		"Please use kettle init --gcp-project-id <id> to specify a project.", projectListLimit))
	// }

	projectIDs := map[string]string{}
	for _, project := range results {
		projectIDs[project.Name] = project.ProjectID
	}
	return projectIDs, nil
}
