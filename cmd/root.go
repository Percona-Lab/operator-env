package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	logVerbose    bool
	nodes         int
	engineVersion string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "openv",
	Short: "Main command for calling subcommands",
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
	rootCmd.PersistentFlags().BoolVarP(&logVerbose,
		"verbose",
		"v",
		false,
		"show debug level logs in developer friendly format",
	)
	upCmd.Flags().IntVar(&nodes,
		"nodes",
		1,
		"set nodes count",
	)
	upCmd.Flags().StringVar(&engineVersion,
		"engv",
		"",
		"set Orchestration platform engine version",
	)
}
