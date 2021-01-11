package cmd

import (
	"fmt"
	"os"

	"github.com/operatorai/operator/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "operator",
	Short: "A CLI tool for creating http functions or services",
	Long: "\n🎯 Operator is a command line tool for creating and deploying" +
		"\n http functions or services.",
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
