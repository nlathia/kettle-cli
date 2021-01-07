package config

// Shared constants
const (
	Version          = "v0.0.2-alpha"
	DeploymentConfig = "operator.config"
)

// Cloud providers
const (
	CloudProvider = "cloud_provider"
	GoogleCloud   = "gcloud"
	AWS           = "aws"
)

var CloudProviderNames = map[string]string{
	"Google Cloud (GCP)":        GoogleCloud,
	"Amazon Web Services (AWS)": AWS,
}

// Deployment types
const (
	DeploymentType      = "deployment_type"
	GoogleCloudFunction = "functions"
	GoogleCloudRun      = "run"
	AWSLambda           = "lambda"
)

var DeploymentNames = map[string]map[string]string{
	GoogleCloud: map[string]string{
		"Google Cloud Function": GoogleCloudFunction,
		"Google Cloud Run":      GoogleCloudRun,
	},
	AWS: map[string]string{
		"AWS Lambda": AWSLambda,
	},
}

// Supported languages
// The values depend on the cloud provider (e.g., "python37" or "python3.7")
const (
	Runtime = "language"
)

var RuntimeNames = map[string]map[string]string{
	GoogleCloudFunction: map[string]string{
		"Python (3.7)": "python37",
		"Go (1.13)":    "go113",
	},
	GoogleCloudRun: map[string]string{
		"Python (3.7)": "python37",
		"Go (1.13)":    "go113",
	},
	AWSLambda: map[string]string{
		"Python (3.7)": "python3.7",
		"Go (1.13)":    "go113",
	},
}

// Cloud-specific config
const (
	// Google Cloud deployments
	ProjectID        = "project_id"
	DeploymentRegion = "region"

	// AWS deployments
	RoleArn = "iam_role_arn"
)
