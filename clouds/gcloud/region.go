package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

func SetDeploymentRegion(settings *config.Settings) error {
	if settings.DeploymentRegion != "" {
		return nil
	}

	regions, err := getGoogleCloudRegions()
	if err != nil {
		return err
	}

	region, err := command.PromptForValue("Deployment region", regions, false)
	if err != nil {
		return err
	}

	settings.DeploymentRegion = region
	return nil
}

// gcloud functions regions list --format="json"
func getGoogleCloudRegions() (map[string]string, error) {
	output, err := command.ExecuteWithResult("gcloud", []string{
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
