package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

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
	configValues, _ = config.ReadSettings()
}

func validateCreateArgs(cmd *cobra.Command, args []string) error {
	if configValues == nil {
		return errors.New("config not found. Please run operator init.")
	}

	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a name")
	}

	// Construct the path where we are going to generate the boiler plate
	var err error // Avoid shadowing global directoryPath
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
	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Set the directory and function name
	configValues.Name = templates.CreateFunctionName(args)
	configValues.FunctionName = templates.CreateEntryFunctionName(args, configValues.Runtime)

	// Print out the config
	fmt.Println("🎇  Type: ", configValues.DeploymentType)
	fmt.Println("🎇  Language: ", configValues.Runtime)

	// Create a directory with the function name
	err := os.Mkdir(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Collect template root path and files
	templateRoot, templateFiles, err := getTemplateFiles()
	if err != nil {
		return err
	}

	for _, assetName := range templateFiles {
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

func getTemplateFiles() (string, []string, error) {
	// Iterate on all of the template files
	// Root: templates/<language>/<cloud-provider>/<type>/
	templateRoot := fmt.Sprintf(
		"templates/%s/%s/%s/",
		strings.Replace(configValues.Runtime, ".", "", 1),
		configValues.CloudProvider,
		configValues.DeploymentType,
	)

	assetNames := templates.AssetNames()
	templateFiles := []string{}
	for _, assetName := range assetNames {
		if strings.Contains(assetName, templateRoot) {
			templateFiles = append(templateFiles, assetName)
		}
	}
	if len(templateFiles) == 0 {
		return "", nil, errors.New(fmt.Sprintf("no matching template for: %s", templateRoot))
	}
	return templateRoot, templateFiles, nil
}

func cleanUp(directoryPath string, err error) error {
	cleanupErr := os.RemoveAll(directoryPath)
	if cleanupErr != nil {
		fmt.Println("\n  Failed to clean up: ", directoryPath, cleanupErr)
	}
	return err
}
