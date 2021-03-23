package clouds

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/operatorai/kettle-cli/clouds/gcloud"
	"github.com/operatorai/kettle-cli/settings"
)

type GoogleCloud struct{}

func (GoogleCloud) GetService(deploymentType string) (Service, error) {
	switch deploymentType {
	case "function":
		return gcloud.GoogleCloudFunction{}, nil
	case "run":
		return gcloud.GoogleCloudRun{}, nil
	}
	return nil, errors.New(fmt.Sprintf("unimplemented service: %s", deploymentType))
}

func (GoogleCloud) Setup(stg *settings.Settings) error {
	_, err := exec.LookPath("gcloud")
	if err != nil {
		return errors.New(fmt.Sprintf("please install the gcloud cli: %s", err))
	}
	if stg.GoogleCloud == nil {
		stg.GoogleCloud = &settings.GoogleCloudSettings{}
	}
	if err := gcloud.SetProjectID(stg.GoogleCloud); err != nil {
		return err
	}
	if err := gcloud.SetDeploymentRegion(stg.GoogleCloud); err != nil {
		return err
	}
	return nil
}
