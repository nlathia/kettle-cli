package clouds

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/preferences"

	"github.com/janeczku/go-spinner"
)

var GcpConfigChoices = []*preferences.ConfigChoice{
	{
		// Pick a Google Cloud Project
		Label:             "Google Cloud Project",
		Key:               config.ProjectID,
		FlagKey:           "gcp-project-id",
		FlagDescription:   "The id of the GCP project to use",
		CollectValuesFunc: getGoogleCloudProjects,
		ValidationFunc:    isActiveGoogleCloudProject,
	},
	{
		// Pick a deployment region
		Label:             "Deployment Region",
		Key:               config.DeploymentRegion,
		FlagKey:           "deployment-region",
		FlagDescription:   "The name of the GCP deployment region to use",
		CollectValuesFunc: getGoogleCloudRegions,
		ValidationFunc:    isValidGoogleCloudRegion,
	},
}

func getGoogleCloudProjects() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud projects...")
	defer s.Stop()

	// gcloud projects list --format="json"
	projectListLimit := 50
	output, err := executeCommandWithResult("gcloud", []string{
		"projects",
		"list",
		"--format=\"json\"",
		fmt.Sprintf("--limit=%d", projectListLimit),
	})
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
	if len(results) >= projectListLimit {
		// Bail if the user has too many projects to 'reasonably' display
		// @TODO add this back in
		return nil, errors.New(fmt.Sprintf("you have %d or more Google Cloud projects. "+
			"Please use operator init --gcp-project-id <id> to specify a project.", projectListLimit))
	}

	projectIDs := map[string]string{}
	for _, project := range results {
		displayName := fmt.Sprintf("%s (%s)", project.Name, project.ProjectID)
		projectIDs[displayName] = project.Name
	}
	return projectIDs, nil
}

func isActiveGoogleCloudProject(projectID string) error {
	s := spinner.StartNew(fmt.Sprintf("Checking Google Cloud project: %s...", projectID))
	defer s.Stop()

	// gcloud projects describe <id> --format="json"
	output, err := executeCommandWithResult("gcloud", []string{
		"projects",
		"describe",
		projectID,
		"--format=\"json\"",
	})
	if err != nil {
		return err
	}

	var result struct {
		LifecycleState string `json:"lifecycleState"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}
	if result.LifecycleState != "ACTIVE" {
		return errors.New("Project is not currently active")
	}
	return nil
}

func getGoogleCloudRegions() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud regions...")
	defer s.Stop()

	// gcloud functions regions list --format="json"
	output, err := executeCommandWithResult("gcloud", []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	})
	if err != nil {
		return nil, err
	}

	var results []struct {
		DisplayName string `json:"displayName"`
		LocationID  string `json:"locationId"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	regions := map[string]string{}
	for _, region := range results {
		displayName := fmt.Sprintf("%s (%s)", region.DisplayName, region.LocationID)
		regions[displayName] = region.LocationID
	}
	return regions, nil
}

func isValidGoogleCloudRegion(locationID string) error {
	s := spinner.StartNew("Collecting Google Cloud regions...")
	defer s.Stop()

	// gcloud functions regions list --format="json"
	output, err := executeCommandWithResult("gcloud", []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	})
	if err != nil {
		return err
	}

	var results []struct {
		LocationID string `json:"locationId"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return err
	}

	for _, region := range results {
		if region.LocationID == locationID {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("unknown region: %s", locationID))
}
