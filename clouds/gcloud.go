package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/clouds/gcloud"
	"github.com/operatorai/operator/command"
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

func (GoogleCloud) Setup(settings *config.Settings) error {
	err := command.Execute("command", []string{
		"-v",
		"gcloud",
	})
	if err != nil {
		return errors.New("please install the gcloud cli")
	}

	if err := gcloud.SetProjectID(settings); err != nil {
		return err
	}
	if err := gcloud.SetDeploymentRegion(settings); err != nil {
		return err
	}
	return nil
}
