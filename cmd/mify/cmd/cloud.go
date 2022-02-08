package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var initCloudCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Mify Cloud user",
	Long: `Initialize Mify Cloud user and config`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CloudInit(appContext, workspacePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to init cloud: %s\n", err)
			os.Exit(2)
		}
	},
}

var updateKubeconfigCmd = &cobra.Command{
	Use:   "update-kubeconfig",
	Short: "Update ~/.kube/config file",
	Long: `Update ~/.kube/config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CloudUpdateKubeconfig(appContext, workspacePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to update kubernetes config: %s\n", err)
			os.Exit(2)
		}
	},
}

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Use Mify Cloud",
	Long:  `Subcommand to access and use Mify Cloud`,
}

func init() {
	cloudCmd.AddCommand(initCloudCmd)
	cloudCmd.AddCommand(updateKubeconfigCmd)
}
