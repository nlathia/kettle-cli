package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/clouds/gcloud"
	"github.com/operatorai/operator/config"
)

type GoogleCloud struct{}

func (GoogleCloud) GetService(deploymentType string) (Service, error) {
	switch deploymentType {
	case config.GoogleCloudFunction:
		return gcloud.GoogleCloudFunction{}, nil
	case config.GoogleCloudRun:
		return gcloud.GoogleCloudRun{}, nil
	}
	return nil, errors.New(fmt.Sprintf("unimplemented service: %s", deploymentType))
}
