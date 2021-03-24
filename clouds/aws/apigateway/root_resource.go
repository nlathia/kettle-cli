package apigateway

import (
	"errors"

	"github.com/operatorai/kettle-cli/settings"
)

func SetRootResourceID(resources []*RestApiResource, stg *settings.Settings) error {
	if stg.AWS.RestApiRootID != "" {
		return nil
	}
	if stg.AWS.RestApiID == "" {
		return errors.New("rest api id not set")
	}

	resource := getResourceWithPath(resources, "")
	if resource == nil {
		return errors.New("did not find root apigateway resource")
	}

	stg.AWS.RestApiRootID = resource.ID
	return nil
}
