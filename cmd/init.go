package cmd

import (
	"fmt"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/clouds"
	"github.com/operatorai/kettle-cli/settings"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up your project / cloud provider settings",
	Long: `🆕 The kettle CLI tool automatically deploys to your cloud provider.
	
	Use this command to (re-)set up the CLI's cloud provider settings.
	`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Pick cloud provider to set up
	cloudProviderName, err := cli.PromptForValue("Cloud to configure", clouds.SupportedClouds(), false)
	if err != nil {
		return formatError(err)
	}

	cloudProvider, err := clouds.GetCloudProvider(cloudProviderName)
	if err != nil {
		return formatError(err)
	}

	// Read the existing settings
	cloudSettings, err := settings.ReadSettings()
	if err != nil {
		return formatError(err)
	}

	// Reset the values
	if err := cloudProvider.Setup(cloudSettings, true); err != nil {
		return formatError(err)
	}

	if cloudSettings.GoogleCloud.ProdProject == nil {
		panic("Prod project is nil")
	}

	// Write them back
	if err := settings.WriteSettings(cloudSettings); err != nil {
		return formatError(err)
	}

	fmt.Println("✅  Settings updated!")
	return nil
}
