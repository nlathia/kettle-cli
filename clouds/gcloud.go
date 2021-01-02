package clouds

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
)

func gcpSetup() error {

	// A configChoice is defined as:
	type configChoice struct {
		// The label, which is shown in the prompt to the end user
		label string

		// The config key: the selection will be stored in viper using this
		key string

		// A function to collect values if the user does not provide one via a flag
		collectValuesFunc func() (map[string]string, error)

		// A function to validate the choice
		validationFunc func(string) error
	}

	configChoices := []configChoice{
		{
			// Pick a Google Cloud Project
			label:             "Google Cloud Project",
			key:               config.ProjectID,
			collectValuesFunc: getGoogleCloudProjects,
			validationFunc:    isActiveGoogleCloudProject,
		},
		{
			// Pick a deployment region
			label:             "Deployment Region",
			key:               config.DeploymentRegion,
			collectValuesFunc: getGoogleCloudRegions,
			validationFunc:    isValidGoogleCloudRegion,
		},
	}

	// Iterate on all of the remaining choices second, since
	// it is slower to collect their values
	for _, choice := range configChoices {
		values, err := choice.collectValuesFunc()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return err
		}
		value, err := getValue(choice.label, values)
		if err != nil {
			return err
		}
		viper.Set(choice.key, value)
	}

	return nil
}

func getGoogleCloudProjects() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud projects...")
	defer s.Stop()

	// gcloud projects list --format="json"
	projectListLimit := 25
	commandArgs := []string{
		"projects",
		"list",
		"--format=\"json\"",
		fmt.Sprintf("--limit=%d", projectListLimit),
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return nil, err
	}

	type gcloudResults struct {
		ProjectID string `json:"projectId"`
		Name      string `json:"name"`
	}
	var results []gcloudResults
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}
	if len(results) == projectListLimit {
		// Bail if the user has too many projects to 'reasonably' display
		// @TODO add this back in
		return nil, errors.New(fmt.Sprintf("you have %d or more Google Cloud projects. "+
			"Please use operator init --project-id <id> to specify a project.", projectListLimit))
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
	commandArgs := []string{
		"projects",
		"describe",
		projectID,
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return err
	}

	type gcloudResult struct {
		LifecycleState string `json:"lifecycleState"`
	}
	var result gcloudResult
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
	commandArgs := []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return nil, err
	}

	type gcloudResults struct {
		DisplayName string `json:"displayName"`
		LocationID  string `json:"locationId"`
	}
	var results []gcloudResults
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
	commandArgs := []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return err
	}

	type gcloudResults struct {
		LocationID string `json:"locationId"`
	}
	var results []gcloudResults
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
