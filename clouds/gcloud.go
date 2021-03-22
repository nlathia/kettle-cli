package clouds

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/operatorai/kettle-cli/clouds/gcloud"
	"github.com/operatorai/kettle-cli/config"
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

func (GoogleCloud) Setup(settings *config.Settings) error {
	_, err := exec.LookPath("gcloud")
	if err != nil {
		return errors.New(fmt.Sprintf("please install the gcloud cli: %s", err))
	}

	if err := gcloud.SetProjectID(settings); err != nil {
		return err
	}
	if err := gcloud.SetDeploymentRegion(settings); err != nil {
		return err
	}
	return nil
}
