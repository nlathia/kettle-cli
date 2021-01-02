package config

import (
	mapset "github.com/deckarep/golang-set"
)

const (
	// viper config keys
	CloudProvider  = "cloud_provider"
	DeploymentType = "deployment_type"
	Runtime        = "language"

	// Google Cloud deployments
	GoogleCloud      = "gcloud"
	ProjectID        = "project_id"
	DeploymentRegion = "region"

	// AWS deployments
	// @TODO

	// Deployment types
	GoogleCloudFunction = "functions"
	GoogleCloudRun      = "run"
	AWSLambda           = "lambda"

	// Supported languages
	Python = "python37"
	GoLang = "go113"

	// Service config file name
	DeploymentConfig = "operator.config"
)

var DeploymentTypes = mapset.NewSetWith(
	GoogleCloudFunction,
	GoogleCloudRun,
	AWSLambda,
)

var DeploymentNames = map[string]string{
	"Google Cloud Function": GoogleCloudFunction,
	"Google Cloud Run":      GoogleCloudRun,
	"AWS Lambda":            AWSLambda,
}

var Runtimes = mapset.NewSetWith(
	Python,
	GoLang,
)

var RuntimeNames = map[string]string{
	"Python (3.7)": Python,
	"Go (1.13)":    GoLang,
}
