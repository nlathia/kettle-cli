package aws

import (
	"encoding/json"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

const (
	operatorUsagePlanName = "operator-apigateway-usage-plan"
)

func getUsagePlans(cfg *config.TemplateConfig) (map[string]string, bool, error) {
	output, err := command.ExecuteWithResult("aws", []string{
		"apigateway",
		"get-usage-plans",
		"--output", "json",
	})
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
			if stage.ID == cfg.RestApiID && stage.Stage == "prod" {
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
