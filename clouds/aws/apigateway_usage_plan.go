package aws

import (
	"encoding/json"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

const (
	operatorUsagePlanName = "operator-apigateway-usage-plan"
)

func getUsagePlans(stg *settings.Settings) (map[string]string, bool, error) {
	output, err := cli.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-usage-plans",
		"--output", "json",
	}, "Collecting available usage plans")
	if err != nil {
		if err.Error() == "exit status 254" {
			return map[string]string{}, false, nil
		}
		return nil, false, err
	}

	var results struct {
		Items []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			ApiStages []struct {
				ID    string `json:"apiId"`
				Stage string `json:"stage"`
			} `json:"apiStages"`
		} `json:"items"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, false, err
	}

	operatorUsagePlanExists := false
	usagePlans := map[string]string{}
	for _, result := range results.Items {
		for _, stage := range result.ApiStages {
			if stage.ID == stg.AWS.RestApiID && stage.Stage == "prod" {
				usagePlans[result.Name] = result.ID
				if result.Name == operatorUsagePlanName {
					operatorUsagePlanExists = true
				}
				break
			}
		}
	}

	return usagePlans, operatorUsagePlanExists, nil
}
