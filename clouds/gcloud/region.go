package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/janeczku/go-spinner"
	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

func setDeploymentRegion(cfg *config.TemplateConfig) error {
	if cfg.DeploymentRegion != "" {
		return nil
	}

	regions, err := getGoogleCloudRegions()
	if err != nil {
		return err
	}

	region, err := command.PromptForValue("Deployment region", regions)
	if err != nil {
		return err
	}

	cfg.DeploymentRegion = region
	return nil
}

// gcloud functions regions list --format="json"
func getGoogleCloudRegions() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud regions...")
	defer s.Stop()

	output, err := command.ExecuteWithResult("gcloud", []string{
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
