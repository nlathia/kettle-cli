package config

type TemplateConfig struct {
	// operator init values
	CloudProvider  string `yaml:"cloud_provider"`
	DeploymentType string `yaml:"deployment_type"`
	Runtime        string `yaml:"runtime"`

	// operator create values
	Name         string `yaml:"name"`
	FunctionName string `yaml:"entrypoint"`

	// operator deploy values
	// Cloud variables
	Deployed         string `yaml:"deployed_utc,omitempty"`
	DeploymentRegion string `yaml:"region,omitempty"`

	// GCP Variables
	ProjectID string `yaml:"project_id,omitempty"`

	// AWS Variables
	AccountID         string `yaml:"account_id,omitempty"`
	RoleArn           string `yaml:"role_arn,omitempty"`
	RestApiID         string `yaml:"rest_api_id,omitempty"`
	RestApiRootID     string `yaml:"rest_api_root_id,omitempty"`
	RestApiResourceID string `yaml:"rest_api_resource_id,omitempty"`
}

// A ConfigChoice is used to enumerate a set of preferences
// that can be selected interactively by the user
type ConfigChoice struct {
	// The label, which is shown in the prompt to the end user
	// The config key: the selection will be stored in viper using this
	Label string
	Key   string

	// Flags so that users can define this choice via an input flag
	// e.g. --cloud <value>
	FlagKey         string
	FlagDescription string
	FlagValue       string

	// A function to collect values if the user does not provide one via a flag
	// A function to validate the choice
	CollectValuesFunc func() (map[string]string, error)
	ValidationFunc    func(string) error
}
