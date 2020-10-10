package config

import (
	mapset "github.com/deckarep/golang-set"
)

const (
	// Deployment types
	DeploymentType = "deployment_type"

	GoogleCloudFunction = "cloud_function"
	GoogleCloudRun      = "cloud_run"

	GoogleCloudFunctionName = "Google Cloud Function"
	GoogleCloudRunName      = "Google Cloud Run"

	// Supported languages
	Runtime = "runtime"
	Python  = "python"
	GoLang  = "go"
)

var DeploymentTypes = mapset.NewSetWith(
	GoogleCloudFunction,
	GoogleCloudRun,
)

var DeploymentNames = map[string]string{
	GoogleCloudFunctionName: GoogleCloudFunction,
	GoogleCloudRunName:      GoogleCloudRun,
}

var Runtimes = mapset.NewSetWith(
	Python,
	// GoLang,
)

var RuntimeNames = map[string]string{
	"Python": Python,
	// "Go":     GoLang,
}
