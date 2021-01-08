package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/config"
)

type Cloud interface {
	GetConfig() *config.TemplateConfig

	Setup() error

	Deploy(directory string, config *config.TemplateConfig) error
}

func GetCloudProvider(cloudType string) (Cloud, error) {
	switch cloudType {
	case config.GoogleCloudFunction:
		return GoogleCloudFunction{}, nil
	case config.GoogleCloudRun:
		return GoogleCloudRun{}, nil
	case config.AWSLambda:
		return AWSLambdaFunction{}, nil
	}
	return nil, errors.New(fmt.Sprintf("Unknown cloud: %s", cloudType))
}
