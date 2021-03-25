package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

func SetDeploymentRegion(stg *settings.GoogleCloudSettings) error {
	if stg.DeploymentRegion != "" {
		return nil
	}

	regions, err := getGoogleCloudRegions()
	if err != nil {
		return err
	}

	region, err := cli.PromptForValue("Deployment region", regions, false)
	if err != nil {
		return err
	}

	stg.DeploymentRegion = region
	return nil
}

// gcloud functions regions list --format="json"
func getGoogleCloudRegions() (map[string]string, error) {
	output, err := cli.ExecuteWithResult("gcloud", []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}, "Collecting gcloud regions")
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
