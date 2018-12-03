package cmd

import (
	"context"
	"fmt"
	"github.com/Percona-Lab/operator-env/config"
	"github.com/Percona-Lab/operator-env/ctrl"
	"github.com/Percona-Lab/operator-env/logger"
	"github.com/spf13/cobra"
	"os"
)

var kubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Brings up/down the Kubernetes cluster",
}

var openshiftCmd = &cobra.Command{
	Use:   "openshift",
	Short: "Brings up/down the OpenShift Cluster",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Brings up the Cluster",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(cmd.CommandPath())

		conf := config.Config{
			LogVerbose: logVerbose,
		}
		log := logger.NewLogger(conf.LogVerbose)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		controller, err := ctrl.NewController(log)
		if err != nil {
			log.Error().Err(err).Msg("controller creation error")
			os.Exit(1)
		}

		var platform ctrl.Platform

		switch cmd.Parent() {
		case openshiftCmd:
			platform = ctrl.OpenShift
		case kubernetesCmd:
			platform = ctrl.Kubernetes
		}

		if err := controller.Up(ctx, platform); err != nil {
			log.Error().Err(err).Msg("cluster creation error")
			os.Exit(1)
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Shut down the Cluster",
	Run: func(cmd *cobra.Command, args []string) {
		//ctx, cancel := context.WithCancel(context.Background())
		//defer cancel()
	},
}

func init() {
	rootCmd.AddCommand(kubernetesCmd, openshiftCmd)
	kubernetesCmd.AddCommand(upCmd, downCmd)
	openshiftCmd.AddCommand(upCmd, downCmd)
}
