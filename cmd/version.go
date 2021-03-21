/*
Copyright © 2020 Neal Lathia <neal.lathia@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/operatorai/kettle/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Installed version of operator",
	Long:  `🔢 Prints and installed version of the operator CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
