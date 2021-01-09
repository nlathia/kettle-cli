package config

const (
	// Shared constants
	Version                 = "v0.0.2-alpha"
	DeploymentConfig        = "operator.config"
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

	// Google Cloud deployments
	ProjectID        = "project_id"
	DeploymentRegion = "region"

	// AWS deployments
	RoleArn             = "iam_role_arn"
	RestApiID           = "rest_api_id"
	RestApiRootResource = "rest_api_root_resource"
)

// Mappings between prompts (shown to the user) and values (stored in config)

var CloudProviderNames = map[string]string{
	"Google Cloud (GCP)":        GoogleCloud,
	"Amazon Web Services (AWS)": AWS,
}

var DeploymentNames = map[string]map[string]string{
	GoogleCloud: map[string]string{
		"Google Cloud Function": GoogleCloudFunction,
		"Google Cloud Run":      GoogleCloudRun,
	},
	AWS: map[string]string{
		"AWS Lambda": AWSLambda,
	},
}

var RuntimeNames = map[string]map[string]string{
	GoogleCloudFunction: map[string]string{
		"Python (3.7)": "python37", // Unlike aws, requires "37"
		"Go (1.13)":    "go113",
	},
	GoogleCloudRun: map[string]string{
		"Python (3.7)": "python37", // Unlike aws, requires "37"
		"Go (1.13)":    "go113",
	},
	AWSLambda: map[string]string{
		"Python (3.7)": "python3.7", // Unlike gcloud, requires the "3.7"
		"Go (1.13)":    "go113",
	},
}
