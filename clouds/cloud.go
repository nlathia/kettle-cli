package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/config"
)

type Cloud interface {
	Build(directory string) error
	Deploy(directory string) error
}

func GetCloudProvider(cloudType string) (Cloud, error) {
	switch cloudType {
	case config.GoogleCloudFunction:
		return GoogleCloudFunction{}, nil
	case config.GoogleCloudRun:
		return GoogleCloudRun{}, nil
	}
	return nil, errors.New(fmt.Sprintf("Unknown cloud: %s", cloudType))
}
