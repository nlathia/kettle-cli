package apigateway

import (
	"encoding/json"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

func GetResources(stg *settings.Settings) ([]*RestApiResource, error) {
	output, err := cli.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-resources",
		"--rest-api-id", stg.AWS.RestApiID,
	}, "Collecting API resources")
	if err != nil {
		return nil, err
	}

	var results struct {
		Items []struct {
			Path            string `json:"path"`
			ID              string `json:"id"`
			ResourceMethods struct {
				POST *struct{} `json:"POST"`
			} `json:"resourceMethods"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	resources := []*RestApiResource{}
	for _, result := range results.Items {
		resources = append(resources, &RestApiResource{
			Path:          result.Path,
			ID:            result.ID,
			HasPostMethod: (result.ResourceMethods.POST != nil),
		})
	}
	return resources, nil
}

func getResourceWithPath(resources []*RestApiResource, pathPart string) *RestApiResource {
	for _, resource := range resources {
		if resource.Path == strings.Join([]string{"/", pathPart}, "") {
			return resource
		}
	}
	return nil
}
