package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/config"
)

type Service interface {
	Deploy(directory string, config *config.TemplateConfig) error
}

type Cloud interface {
	Setup(settings *config.Settings) error

	GetService(deploymentType string) (Service, error)
}

func GetCloudProvider(cloudType string) (Cloud, error) {
	switch cloudType {
	case config.GoogleCloud:
		return GoogleCloud{}, nil
	case config.AWS:
		return AmazonWebServices{}, nil
	}
	return nil, errors.New(fmt.Sprintf("unimplemented cloud: %s", cloudType))
}
