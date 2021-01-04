package config

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type TemplateConfig struct {
	CloudProvider string `yaml:"cloud_provider"`
	Name          string `yaml:"name"`
	FunctionName  string `yaml:"entrypoint"`
	Runtime       string `yaml:"runtime"`
	Type          string `yaml:"type"`

	// GCP Variables
	ProjectID        string `yaml:"project_id,omitempty"`
	DeploymentRegion string `yaml:"region,omitempty"`
	Deployed         string `yaml:"deployed_utc,omitempty"`
	PackageName      string `yaml:"package,omitempty"`

	// AWS Variables
	IAMRole string `yaml:"iam_role,omitempty"`
}

func GetConfigFilePath(directoryPath string) string {
	return path.Join(directoryPath, DeploymentConfig)
}

func WriteConfig(config *TemplateConfig, directoryPath string) error {
	data, err := yaml.Marshal(config)
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
