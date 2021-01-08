package clouds

import (
	"errors"
	"fmt"

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

func (AmazonWebServices) AddConfig(cfg *config.TemplateConfig) error {
	return nil
}
