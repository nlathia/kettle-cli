/*
Copyright © 2020 Neal Lathia <neal.lathia@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version = "v0.0.23"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Installed version of kettle",
	Long:  `🔢 Prints the installed version of the kettle CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
