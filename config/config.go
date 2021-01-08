package config

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type TemplateConfig struct {
	// Template Variables
	Name           string `yaml:"name"`
	FunctionName   string `yaml:"entrypoint"`
	Runtime        string `yaml:"runtime"`
	DeploymentType string `yaml:"deployment_type"`
	Deployed       string `yaml:"deployed_utc,omitempty"`
	PackageName    string `yaml:"package,omitempty"`

	// Cloud variables
	CloudProvider    string `yaml:"cloud_provider"`
	DeploymentRegion string `yaml:"region,omitempty"`

	// GCP Variables
	ProjectID string `yaml:"project_id,omitempty"`

	// AWS Variables
	RoleArn   string `yaml:"role_arn,omitempty"`
	RestApiID string `yaml:"rest_api_id,omitempty"`
}

func GetConfigFilePath(directoryPath string) string {
	return path.Join(directoryPath, DeploymentConfig)
}

func WriteConfig(cfg *TemplateConfig, directoryPath string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath := GetConfigFilePath(directoryPath)
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig(directoryPath string) (*TemplateConfig, error) {
	filePath := GetConfigFilePath(directoryPath)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	values := TemplateConfig{}
	if err := yaml.Unmarshal(contents, &values); err != nil {
		return nil, err
	}
	return &values, nil
}
