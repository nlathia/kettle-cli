package gcloud

import (
	"encoding/json"
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/settings"
)

func SetProjects(sts *settings.GoogleCloudSettings, overwrite bool) error {
	environments := map[string]*settings.GoogleCloudProject{
		"development": sts.DevProject,
		"production":  sts.ProdProject,
	}

	if !overwrite {
		valuesSet := true
		for _, environment := range environments {
			if environment == nil {
				valuesSet = false
				break
			}
			if environment.ProjectName == "" || environment.ProjectID == "" {
				valuesSet = false
				break
			}
		}
		if valuesSet {
			return nil
		}
	}

	projects, err := getGoogleCloudProjects()
	if err != nil {
		return err
	}

	regions, err := getGoogleCloudRegions()
	if err != nil {
		return err
	}

	sts.ProdProject, err = setupEnvironment("production", projects, regions)
	if err != nil {
		return err
	}

	sts.DevProject, err = setupEnvironment("development", projects, regions)
	if err != nil {
		return err
	}

	fmt.Println(sts)
	return nil
}

func setupEnvironment(name string, projects, regions map[string]string) (*settings.GoogleCloudProject, error) {
	fmt.Printf("\n🔎 Set up a Google Cloud environment: %s\n", name)
	prompt := fmt.Sprintf("Select your project for \"%s\"", name)
	projectName, projectID, err := cli.PromptForKeyValue(prompt, projects)
	if err != nil {
		return nil, err
	}

	prompt = fmt.Sprintf("Select your deployment region for \"%s\"", name)
	region, err := cli.PromptForValue(prompt, regions, false)
	if err != nil {
		return nil, err
	}

	return &settings.GoogleCloudProject{
		ProjectName:      projectName,
		ProjectID:        projectID,
		DeploymentRegion: region,
	}, nil
}

func getGoogleCloudProjects() (map[string]string, error) {
	// gcloud projects list --format="json"
	output, err := cli.ExecuteWithResult("gcloud", []string{
		"projects",
		"list",
		"--format=\"json\"",
		fmt.Sprintf("--limit=%d", 50),
	}, "Collecting gcloud projects")
	if err != nil {
		return nil, err
	}

	var results []struct {
		ProjectID string `json:"projectId"`
		Name      string `json:"name"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	// @TODO re-enable this
	// if len(results) >= projectListLimit {
	// 	// Bail if the user has too many projects to 'reasonably' display
	// 	return nil, errors.New(fmt.Sprintf("you have %d or more Google Cloud projects. "+
	// 		"Please use kettle init --gcp-project-id <id> to specify a project.", projectListLimit))
	// }

	projectIDs := map[string]string{}
	for _, project := range results {
		projectIDs[project.Name] = project.ProjectID
	}
	return projectIDs, nil
}

// gcloud functions regions list --format="json"
func getGoogleCloudRegions() (map[string]string, error) {
	output, err := cli.ExecuteWithResult("gcloud", []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}, "Collecting gcloud regions")
	if err != nil {
		return nil, err
	}

	var results []struct {
		DisplayName string `json:"displayName"`
		LocationID  string `json:"locationId"`
	}
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	regions := map[string]string{}
	for _, region := range results {
		displayName := fmt.Sprintf("%s (%s)", region.DisplayName, region.LocationID)
		regions[displayName] = region.LocationID
	}
	return regions, nil
}

func getEnvironment(stg *settings.Settings, env string) *settings.GoogleCloudProject {
	if env == "prod" || env == "production" {
		return stg.GoogleCloud.ProdProject
	}
	return stg.GoogleCloud.DevProject
}
