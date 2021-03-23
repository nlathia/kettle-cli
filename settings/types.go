package settings

const (
	settingsFileName = "~/.kettle.yaml"
)

// Debug mode (kettle <command> --debug)
var DebugMode bool

// Settings are values that do not change (often) across multiple deployments

type GoogleCloudSettings struct {
	ProjectName      string `yaml:"project_name,omitempty"`
	ProjectID        string `yaml:"project_id,omitempty"`
	DeploymentRegion string `yaml:"region,omitempty"`
}

type AWSSettings struct {
	AccountID     string `yaml:"account_id,omitempty"`
	RoleArn       string `yaml:"role_arn,omitempty"`
	RestApiID     string `yaml:"rest_api_id,omitempty"`
	RestApiRootID string `yaml:"rest_api_root_id,omitempty"`
}

type Settings struct {
	GoogleCloud *GoogleCloudSettings `yaml:"gcloud,omitempty"`
	AWS         *AWSSettings         `yaml:"aws,omitempty"`
}

// Values that are specific to each deployment
// type TemplateConfig struct {

// 	Settings *Settings `yaml:"settings"`

// 	// template create values
// 	Name         string `yaml:"name"`
// 	FunctionName string `yaml:"entrypoint"`

// 	// AWS variables
//
// 	RestApiResourceID string `yaml:"rest_api_resource_id,omitempty"`
// }
