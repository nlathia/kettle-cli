package clouds

import (
	"errors"
	"fmt"

	"github.com/operatorai/operator/clouds/aws"
	"github.com/operatorai/operator/command"
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

func (AmazonWebServices) Setup() error {
	err := command.Execute("command", []string{
		"-v",
		"aws",
	})
	if err != nil {
		return errors.New("please install the aws cli")
	}
	return nil
}
