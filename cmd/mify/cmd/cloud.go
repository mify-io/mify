package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	deployEnv   string
	confEnv   string
	shellEnv   string
	forwardProxy string
	listenPort string
)

var initCloudCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Mify Cloud user",
	Long:  `Initialize Mify Cloud user and config`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CloudInit(appContext); err != nil {
			fmt.Fprintf(os.Stderr, "failed to init cloud: %s\n", err)
			os.Exit(2)
		}
	},
}

var updateKubeconfigCmd = &cobra.Command{
	Use:   "update-kubeconfig",
	Short: "Update ~/.kube/config file",
	Long:  `Update ~/.kube/config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CloudUpdateKubeconfig(appContext, confEnv); err != nil {
			fmt.Fprintf(os.Stderr, "failed to update kubernetes config: %s\n", err)
			os.Exit(2)
		}
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy [service]...",
	Short: "Deploy code to cloud",
	Long:  `Deploy code to cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.DeployMany(appContext, deployEnv, args); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run deploy: %s\n", err)
			os.Exit(2)
		}
	},
}

var nsShellCmd = &cobra.Command{
	Use:   "ns-shell",
	Short: "Run shell in cloud namespace",
	Long:  `Run shell in cloud namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.NsShell(appContext, shellEnv, forwardProxy, listenPort); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start shell: %s\n", err)
			os.Exit(2)
		}
	},
}

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Use Mify Cloud",
	Long:  `Subcommand to access and use Mify Cloud`,
	PersistentPreRun: func(*cobra.Command, []string) {
		err := appContext.LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init workspace: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
	deployCmd.PersistentFlags().StringVarP(&deployEnv, "environment", "e", "stage", "Target environment name")
	updateKubeconfigCmd.PersistentFlags().StringVarP(&confEnv, "environment", "e", "stage", "Target environment name")

	cloudCmd.AddCommand(initCloudCmd)
	cloudCmd.AddCommand(updateKubeconfigCmd)
	cloudCmd.AddCommand(deployCmd)

	cloudCmd.AddCommand(nsShellCmd)
	nsShellCmd.PersistentFlags().StringVarP(&shellEnv, "environment", "e", "stage", "Target environment name")
	nsShellCmd.PersistentFlags().StringVarP(&forwardProxy,
		"forward-proxy", "L", "",
		"Proxy remote address from pod, usage: bind-port:remote-host:remote-port")
	nsShellCmd.PersistentFlags().StringVarP(&listenPort,
		"ssh-port", "P", "49222",
		"port for connecting to pod via ssh, default: 49222")
}
