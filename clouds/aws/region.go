package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
	"github.com/spf13/viper"
)

func setDeploymentRegion(cfg *config.TemplateConfig) error {
	if cfg.DeploymentRegion != "" {
		return nil
	}

	regions, err := getAWSRegions()
	if err != nil {
		return err
	}

	region, err := command.PromptForValue("Deployment region", regions, false)
	if err != nil {
		return err
	}

	cfg.DeploymentRegion = region
	viper.Set(config.DeploymentRegion, region)
	return nil
}

// aws ec2 describe-regions --output json
func getAWSRegions() (map[string]string, error) {
	output, err := command.ExecuteWithResult("aws", []string{
		"ec2",
		"describe-regions",
		"--output", "json",
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Regions []struct {
			RegionName string `json:"RegionName"`
		} `json:"Regions"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	regions := map[string]string{}
	for _, region := range result.Regions {
		regions[region.RegionName] = region.RegionName
	}
	return regions, nil
}
