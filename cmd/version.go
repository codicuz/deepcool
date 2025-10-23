package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)
var version = "2.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deep Cool version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
