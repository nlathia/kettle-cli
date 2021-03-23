package config

const (
	configFileName = "kettle.json"
)

type Config struct {
	ProjectName string `json:"name"`
	Config      struct {
		Runtime        string `json:"runtime"`
		CloudProvider  string `json:"cloud_provider"`
		DeploymentType string `json:"deployment_type"`
		FunctionName   string `json:"entry_function"`
	} `json:"config"`
	Template []struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	} `json:"template,omitempty"`
}
