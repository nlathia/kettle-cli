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
const (
	Runtime = "language"
	Python  = "python37"
	GoLang  = "go113"
)

var RuntimeNames = map[string]string{
	"Python (3.7)": Python,
	"Go (1.13)":    GoLang,
}

// Google Cloud deployments
const (
	ProjectID        = "project_id"
	DeploymentRegion = "region"

	// AWS deployments
	IAMRole = "iam_role"
)
