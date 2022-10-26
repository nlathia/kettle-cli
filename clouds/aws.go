package clouds

import (
	"fmt"
	"os/exec"

	"github.com/operatorai/kettle-cli/clouds/aws"
	"github.com/operatorai/kettle-cli/settings"
)

type AmazonWebServices struct{}

func (AmazonWebServices) GetService(deploymentType string) (Service, error) {
	switch deploymentType {
	case "lambda":
		return aws.AWSLambdaFunction{}, nil
	}
	return nil, fmt.Errorf("unimplemented service: %s", deploymentType)
}

func (AmazonWebServices) Setup(stg *settings.Settings, overwrite bool) error {
	_, err := exec.LookPath("aws")
	if err != nil {
		return fmt.Errorf("please install the aws cli: %s", err)
	}
	if stg.AWS == nil {
		stg.AWS = &settings.AWSSettings{}
	}
	if err := aws.SetAccountID(stg.AWS, overwrite); err != nil {
		return err
	}
	if err := aws.SetDeploymentRegion(stg.AWS, overwrite); err != nil {
		return err
	}
	return nil
}
