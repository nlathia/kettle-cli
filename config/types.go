package config

const (
	configFileName = "kettle.json"
)

// Config are values that are specific to individual projects
// and are therefore stored in a config file, one per project

type Config struct {
	ProjectName string `json:"name"`
	Config      struct {
		Runtime        string `json:"runtime"`
		CloudProvider  string `json:"cloud_provider"`
		DeploymentType string `json:"deployment_type"`
		AWS            struct {
			RestApiResourceID string `json:"rest_api_resource_id,omitempty"`
		} `json:"aws,omitempty"`
	} `json:"config"`
	Template []struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Key    string `json:"key"`
		Value  string `json:"value"`
		Style  string `json:"format,omitempty"`
	} `json:"template,omitempty"`
}
