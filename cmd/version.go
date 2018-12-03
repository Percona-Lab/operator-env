package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represent command for showing application version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
