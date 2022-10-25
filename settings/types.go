package settings

// Debug mode (kettle <command> --debug)
var DebugMode bool

// Settings are values that do not change across multiple deployments
// and are therefore stored in a settings file

type GoogleCloudProject struct {
	ProjectName string `yaml:"project_name,omitempty"`
	ProjectID   string `yaml:"project_id,omitempty"`
}

type GoogleCloudSettings struct {
	DevProject       GoogleCloudProject `yaml:"dev_environment,omitempty"`
	ProdProject      GoogleCloudProject `yaml:"prod_environment,omitempty"`
	DeploymentRegion string             `yaml:"region,omitempty"`
}

type AWSSettings struct {
	AccountID        string `yaml:"account_id,omitempty"`
	RoleArn          string `yaml:"role_arn,omitempty"`
	RestApiID        string `yaml:"rest_api_id,omitempty"`
	RestApiRootID    string `yaml:"rest_api_root_id,omitempty"`
	DeploymentRegion string `yaml:"region,omitempty"`
}

type Settings struct {
	GoogleCloud *GoogleCloudSettings `yaml:"gcloud,omitempty"`
	AWS         *AWSSettings         `yaml:"aws,omitempty"`
}
