package clouds

import (
	"fmt"

	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
)

type Service interface {
	Deploy(directory string, cfg *config.Config, stg *settings.Settings, env string) error
}

type Cloud interface {
	Setup(settings *settings.Settings, overwrite bool) error

	GetService(deploymentType string) (Service, error)
}

func GetCloudProvider(cloudType string) (Cloud, error) {
	switch cloudType {
	case "gcloud":
		return GoogleCloud{}, nil
	case "aws":
		return AmazonWebServices{}, nil
	}
	return nil, fmt.Errorf("unimplemented cloud: %s", cloudType)
}

func SupportedClouds() map[string]string {
	return map[string]string{
		"Google Cloud":        "gcloud",
		"Amazon Web Services": "aws",
	}
}
