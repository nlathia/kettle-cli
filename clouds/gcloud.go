package clouds

import (
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
	return nil, fmt.Errorf("unimplemented service: %s", deploymentType)
}

func (GoogleCloud) Setup(stg *settings.Settings, overwrite bool) error {
	_, err := exec.LookPath("gcloud")
	if err != nil {
		return fmt.Errorf("please install the gcloud cli: %s", err)
	}
	if stg.GoogleCloud == nil {
		stg.GoogleCloud = &settings.GoogleCloudSettings{}
	}
	if err := gcloud.SetProjects(stg.GoogleCloud, overwrite); err != nil {
		return err
	}
	return nil
}
