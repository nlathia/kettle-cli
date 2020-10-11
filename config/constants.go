package config

import (
	mapset "github.com/deckarep/golang-set"
)

const (
	// Cloud providers
	CloudProvider = "cloud_provider"
	GoogleCloud   = "gcloud"

	// Deployment types
	DeploymentType = "deployment_type"

	GoogleCloudFunction = "cloud_function"
	GoogleCloudRun      = "cloud_run"

	// Supported languages (just Python right now)
	Runtime = "runtime"
	Python  = "python"

	// Service config file name
	DeploymentConfig = "operator.config"
)

var DeploymentTypes = mapset.NewSetWith(
	GoogleCloudFunction,
	GoogleCloudRun,
)

var DeploymentNames = map[string]string{
	"Google Cloud Function": GoogleCloudFunction,
	"Google Cloud Run":      GoogleCloudRun,
}

var CloudProviders = map[string]string{
	GoogleCloudFunction: GoogleCloud,
	GoogleCloudRun:      GoogleCloud,
}

var Runtimes = mapset.NewSetWith(
	Python,
)

var RuntimeNames = map[string]string{
	"Python": Python,
}
