package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/operatorai/kettle-cli/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kettle",
	Short: "A CLI tool for creating http functions or services",
	Long: "\n🎯 The kettle CLI creates machine learning pipelines" +
		"\n or microservices from templates.",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&config.DebugMode, "debug", false, "Enable debug mode")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.operator.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func formatError(err error) error {
	fmt.Println(fmt.Sprintf("\n❌ %s", err.Error()))
	return nil
}
