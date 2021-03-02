package clouds

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/operatorai/operator/clouds/aws"
	"github.com/operatorai/operator/config"
)

type AmazonWebServices struct{}

func (AmazonWebServices) GetService(deploymentType string) (Service, error) {
	switch deploymentType {
	case config.AWSLambda:
		return aws.AWSLambdaFunction{}, nil
	}
	return nil, errors.New(fmt.Sprintf("unimplemented service: %s", deploymentType))
}

func (AmazonWebServices) Setup(settings *config.Settings) error {
	_, err := exec.LookPath("aws")
	if err != nil {
		return errors.New("please install the gcloud cli")
	}
	if err != nil {
		return errors.New("please install the aws cli")
	}

	if err := aws.SetAccountID(settings); err != nil {
		return err
	}
	if err := aws.SetDeploymentRegion(settings); err != nil {
		return err
	}
	return nil
}
