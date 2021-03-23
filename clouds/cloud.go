package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/templates"
)

type Service interface {
	Deploy(directory string, config *templates.Template) error
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
