package config

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
