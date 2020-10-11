package config

import (
	mapset "github.com/deckarep/golang-set"
)

const (
	// viper config keys
	CloudProvider  = "cloud_provider"
	DeploymentType = "deployment_type"
	Runtime        = "runtime"

	// Google Cloud deployments
	GoogleCloud = "gcloud"
	ProjectID   = "project_id"

	// Deployment types
	GoogleCloudFunction = "cloud_function"
	GoogleCloudRun      = "cloud_run"

	// Supported languages (just Python right now)
	Python = "python37"

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
	"Python (3.7)": Python,
}
