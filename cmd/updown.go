package cmd

import (
	"context"
	"github.com/Percona-Lab/operator-env/config"
	"github.com/Percona-Lab/operator-env/ctrl"
	"github.com/Percona-Lab/operator-env/logger"
	"github.com/spf13/cobra"
	"os"
)

// upCmd
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Brings up the Cluster",
	Long:  ``,
}

// kubernetesCmd
var kubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Brings up/down the Kubernetes cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{
			LogVerbose: logVerbose,
		}
		log := logger.NewLogger(cfg.LogVerbose)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		controller, err := ctrl.NewController(log, cfg)
		if err != nil {
			log.Error().Err(err).Msg("controller creation error")
			os.Exit(1)
		}

		if err := controller.Up(ctx, ctrl.Kubernetes); err != nil {
			log.Error().Err(err).Msg("cluster creation error")
			os.Exit(1)
		}
	},
}

// openshiftCmd
var openshiftCmd = &cobra.Command{
	Use:   "openshift",
	Short: "Brings up/down the OpenShift Cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{
			LogVerbose: logVerbose,
		}
		log := logger.NewLogger(cfg.LogVerbose)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		controller, err := ctrl.NewController(log, cfg)
		if err != nil {
			log.Error().Err(err).Msg("controller creation error")
			os.Exit(1)
		}

		if err := controller.Up(ctx, ctrl.OpenShift); err != nil {
			log.Error().Err(err).Msg("cluster creation error")
			os.Exit(1)
		}
	},
}

// downCmd
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Shut down the Cluster",
	Run: func(cmd *cobra.Command, args []string) {
		//ctx, cancel := context.WithCancel(context.Background())
		//defer cancel()
	},
}

func init() {
	rootCmd.AddCommand(upCmd, downCmd)
	upCmd.AddCommand(kubernetesCmd, openshiftCmd)
}
