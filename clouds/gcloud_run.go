package clouds

import "github.com/operatorai/operator/config"

type GoogleCloudRun struct{}

func (GoogleCloudRun) Build(directory string, config *config.TemplateConfig) error {
	return nil
}

func (GoogleCloudRun) Deploy(directory string, config *config.TemplateConfig) error {
	return nil
}
