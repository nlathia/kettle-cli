package config

// Values that do not change across multiple deployments
type Settings struct {
	CloudProvider    string `yaml:"cloud_provider"`
	DeploymentType   string `yaml:"deployment_type"`
	Runtime          string `yaml:"runtime"`
	DeploymentRegion string `yaml:"region"`

	// GCP Variables
	ProjectID string `yaml:"project_id,omitempty"`

	// AWS Variables
	AccountID     string `yaml:"account_id,omitempty"`
	RoleArn       string `yaml:"role_arn,omitempty"`
	RestApiID     string `yaml:"rest_api_id,omitempty"`
	RestApiRootID string `yaml:"rest_api_root_id,omitempty"`
}

// Values that are specific to each deployment
type TemplateConfig struct {
	Settings *Settings `yaml:"settings"`

	// operator create values
	Name         string `yaml:"name"`
	FunctionName string `yaml:"entrypoint"`

	// AWS variables
	Deployed          string `yaml:"deployed_utc,omitempty"`
	RestApiResourceID string `yaml:"rest_api_resource_id,omitempty"`
}
