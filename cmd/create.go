package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/templates"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new directory with boiler plate code to deploy.",
	Long: `The operator CLI tool can automatically create a directory
 with all of the boiler plate that you need to get started.
	
The create command will create a directory with all the code to get you started.`,
	Args: validateCreateArgs,
	RunE: runCreate,
}

// When we create a deployment, we store everything in a yaml config file
// we will need this later to deploy the function
var configValues *config.TemplateConfig
var directoryPath string

func init() {
	rootCmd.AddCommand(createCmd)
	err := config.Read()
	if err != nil {
		// User has not run operator init
		return
	}

	// Set up the config for this template
	configValues = &config.TemplateConfig{}
	createCmd.Flags().StringVar(&configValues.CloudProvider, "cloud", viper.GetString(config.CloudProvider), "The name of the cloud provider")
	createCmd.Flags().StringVar(&configValues.Type, "type", viper.GetString(config.DeploymentType), "The type of deployment to create")
	createCmd.Flags().StringVar(&configValues.Runtime, "runtime", viper.GetString(config.Runtime), "The function's runtime language")

	// Google Cloud specific flags
	createCmd.Flags().StringVar(&configValues.ProjectID, "project-id", viper.GetString(config.ProjectID), "The gcloud project use")
	createCmd.Flags().StringVar(&configValues.DeploymentRegion, "region", viper.GetString(config.DeploymentRegion), "The region to deploy to")
}

func validateCreateArgs(cmd *cobra.Command, args []string) error {
	if configValues == nil {
		return errors.New("please run operator init to set up this tool")
	}

	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a name")
	}

	// Construct the path where we are going to generate the boiler plate
	var err error
	directoryPath, err = templates.GetRelativeDirectory(args[0])
	if err != nil {
		return err
	}

	// Validate that the function path does *not* already exist
	exists, err := templates.PathExists(directoryPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("directory already exists")
	}

	// Validate the input cloud provider
	if err := mapContainsValue(configValues.CloudProvider, config.CloudProviderNames); err != nil {
		return err
	}

	// Validate the selected runtime is supported
	if err := mapContainsValue(configValues.Runtime, config.RuntimeNames); err != nil {
		return err
	}

	// Validate the selected type of deployment
	if err := mapContainsValue(configValues.Type, config.DeploymentNames[configValues.CloudProvider]); err != nil {
		return err
	}

	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Set the directory and function name
	configValues.Name = templates.CreateFunctionName(args)
	configValues.FunctionName = templates.CreateEntryFunctionName(args, configValues.Runtime)

	// Print out the config
	fmt.Println("🎇  Type: ", configValues.Type)
	fmt.Println("🎇  Language: ", configValues.Runtime)

	// Create a directory with the function name
	err := os.Mkdir(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Iterate on all of the template files
	// Root: templates/<language>/<cloud-provider>/<type>/
	templateRoot := fmt.Sprintf(
		"templates/%s/%s/%s",
		configValues.Runtime,
		configValues.CloudProvider,
		configValues.Type,
	)
	assetNames := templates.AssetNames()

	// Basic error checking
	numAssetFiles := 0
	for _, assetName := range assetNames {
		if strings.Contains(assetName, templateRoot) {
			numAssetFiles++
		}
	}
	if numAssetFiles == 0 {
		return errors.New(fmt.Sprintf("no template for: %s", templateRoot))
	}

	for _, assetName := range assetNames {
		// Skip assets that are not part of the desired template
		if !strings.Contains(assetName, templateRoot) {
			continue
		}

		// Create the target path
		targetPath := strings.Replace(assetName, templateRoot, "", 1)
		targetPath = path.Join(directoryPath, targetPath)

		// Create the parent directory
		parentDir, _ := path.Split(targetPath)
		err = os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return cleanUp(directoryPath, err)
		}

		// Read the asset out of go-bindata
		content, err := templates.Asset(assetName)
		if err != nil {
			return cleanUp(directoryPath, err)
		}

		// Create the file itself
		f, err := os.Create(targetPath)
		if err != nil {
			return cleanUp(directoryPath, err)
		}
		defer f.Close()

		// Render the template into the target file
		tmpl, err := template.New(assetName).Parse(string(content))
		if err != nil {
			return cleanUp(directoryPath, err)
		}

		err = tmpl.Execute(f, configValues)
		if err != nil {
			return cleanUp(directoryPath, err)
		}

		// If it is a .sh file, chmod u+x it
		if strings.HasSuffix(targetPath, ".sh") {
			if err := os.Chmod(targetPath, 0700); err != nil {
				return cleanUp(directoryPath, err)
			}
		}
	}

	err = config.WriteConfig(configValues, directoryPath)
	if err != nil {
		return cleanUp(directoryPath, err)
	}
	fmt.Println("\n✅  Created: ", directoryPath)
	return nil
}

func cleanUp(directoryPath string, err error) error {
	cleanupErr := os.RemoveAll(directoryPath)
	if cleanupErr != nil {
		fmt.Println("\n⚠️  Failed to clean up: ", directoryPath, cleanupErr)
	}
	return err
}
