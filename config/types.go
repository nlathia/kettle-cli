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
		PythonManager  string `json:"python_manager,omitempty"`
		CloudProvider  string `json:"cloud_provider"`
		DeploymentType string `json:"deployment_type"`
		EntryFunction  string `json:"entry_function"`
		AWS            struct {
			RestApiResourceID string `json:"rest_api_resource_id,omitempty"`
		} `json:"deploy_settings,omitempty"`
	} `json:"config"`
	Template []struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Key    string `json:"key"`
		Value  string `json:"value"`
		Style  string `json:"format,omitempty"`
	} `json:"template,omitempty"`
}
