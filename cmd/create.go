package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/templates"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new directory with boiler plate.",
	Long: `The operator CLI tool can automatically create a directory
 with all of the boiler plate that you need to get started.
	
	The create command will create a directory with all the code to get you started.`,
	Args: createArgs,
	RunE: runCreate,
}

var runTimeLanguage string
var serviceType string

func init() {
	rootCmd.AddCommand(createCmd)
	config.Read()

	// Flag to flip between golang and python runtime
	createCmd.Flags().StringVar(&runTimeLanguage, "runtime", viper.GetString(config.Runtime), "The function's runtime language")
	createCmd.Flags().StringVar(&serviceType, "type", viper.GetString(config.DeploymentType), "The type of deployment to create")
}

func createArgs(cmd *cobra.Command, args []string) error {
	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a name")
	}

	// Validate that the function path does *not* already exist
	exists, err := getDirectoryExists(args)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("directory already exists")
	}

	// Validate the selected runtime is supported
	exists = config.Runtimes.Contains(runTimeLanguage)
	if !exists {
		return fmt.Errorf("runtime (%v) needs to be one of (%v)", runTimeLanguage, config.Runtimes.ToSlice())
	}

	// Validate the selected type of deployment
	exists = config.DeploymentTypes.Contains(serviceType)
	if !exists {
		return fmt.Errorf("type (%v) needs to be one of (%v)", serviceType, config.DeploymentTypes.ToSlice())
	}
	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Print out the config
	fmt.Println("🎇  Type: ", serviceType)
	fmt.Println("🎇  Language: ", runTimeLanguage)

	// Populate the template values
	configValues := &config.TemplateConfig{
		CloudName:     "gcloud",
		DirectoryName: getFunctionName(args),
		FunctionName:  getEntryFunctionName(args, runTimeLanguage),
		Runtime:       runTimeLanguage,
		Type:          serviceType,
	}
	if configValues.Runtime == config.GoLang {
		configValues.PackageName = strings.ToLower(
			strcase.ToLowerCamel(configValues.FunctionName),
		)
	}

	// Create a directory with the function name
	functionPath, err := getDirectoryPath(args)
	if err != nil {
		return err
	}

	err = os.Mkdir(functionPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Iterate on all of the template files
	templateRoot := fmt.Sprintf(
		"templates/gcloud/%s/",
		serviceType,
	)
	assetNames := templates.AssetNames()
	for _, assetName := range assetNames {

		// Skip assets that are not part of the desired template
		if !strings.Contains(assetName, templateRoot) {
			continue
		}

		// Create the target path
		targetPath := strings.Replace(assetName, templateRoot, "", 1)
		targetPath = path.Join(functionPath, targetPath)

		// Create the parent directory
		parentDir, _ := path.Split(targetPath)
		exists, err := pathExists(parentDir)
		if err != nil {
			return err
		}
		if !exists {
			fmt.Println("🎯  Creating: ", parentDir)
		}

		err = os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Read the asset out of go-bindata
		content, err := templates.Asset(assetName)
		if err != nil {
			return err
		}

		// Render the template into the target file
		tmpl, err := template.New(assetName).Parse(string(content))
		if err != nil {
			return err
		}

		err = tmpl.Execute(f, configValues)
		if err != nil {
			return err
		}

		// If it is a .sh file, chmod u+x it
		if strings.HasSuffix(targetPath, ".sh") {
			if err := os.Chmod(targetPath, 0700); err != nil {
				return err
			}
		}
	}

	err = config.WriteConfig(configValues, functionPath)
	if err != nil {
		return err
	}
	fmt.Println("\n✅  Created: ", functionPath)
	return nil
}
