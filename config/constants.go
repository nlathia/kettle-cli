package config

const (
	// Shared constants
	Version                 = "v0.0.3"
	DeploymentConfig        = "operator.json"
	PromptNoneOfTheseOption = "None of these (create a new one)"

	// Cloud providers
	CloudProvider = "cloud_provider"
	GoogleCloud   = "gcloud"
	AWS           = "aws"

	// Deployment details
	Runtime = "language"

	DeploymentType      = "deployment_type"
	GoogleCloudFunction = "function"
	GoogleCloudRun      = "run"
	AWSLambda           = "lambda"

	// Cloud variables
	DeploymentRegion = "region"

	// Google Cloud deployments
	ProjectID = "project_id"

	// AWS deployments
	AccountID           = "account_id"
	RoleArn             = "role_arn"
	RestApiID           = "rest_api_id"
	RestApiRootResource = "rest_api_root_resource"
)

// Debug mode
var DebugMode bool
