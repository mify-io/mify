package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate <service-name> <up [N]|down [N]> [-- -help -extra -args]",
	Short: "Run database migration on service locally",
	Long: `Run database migration on service locally`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			_ = cmd.Help()
			fmt.Fprintf(os.Stderr, "error: missing service-name and command\n")
			os.Exit(1)
		}
		serviceName := args[0]
		command := args[1]
		extraArgs := args[2:]
		if err := mify.ToolMigrate(appContext, workspacePath, serviceName, command, extraArgs); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run migrate: %s\n", err)
			os.Exit(2)
		}
	},
}

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Run external tool needed for development",
	Long:  `Run some external tool which is needed for development`,
}

func init() {
	toolCmd.AddCommand(migrateCmd)
}
