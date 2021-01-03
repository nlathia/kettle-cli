package clouds

import "github.com/operatorai/operator/config"

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Setup() error {
	return nil
}

func (AWSLambdaFunction) Deploy(directory string, config *config.TemplateConfig) error {
	return nil
}
