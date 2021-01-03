package clouds

import "github.com/operatorai/operator/config"

type AWSLambdaFunction struct{}

func (AWSLambdaFunction) Setup() error {
	// @TODO: enable selecting whether to create .zip or image-based lambdas
	return nil
}

func (AWSLambdaFunction) Deploy(directory string, config *config.TemplateConfig) error {
	/*
		1. Build the docker image
		2. Push the image to ECR
		3. Deploy the function
	*/

	/*
		Example:
		$ aws lambda create-function \
			--function-name my-function \
			--package-type Image \
			--code # URI of a container image in the Amazon ECR registry.
			--handler index.handler \
			--runtime nodejs12.x \x
			--role arn:aws:iam::123456789012:role/lambda-ex
	*/
	return nil
}
