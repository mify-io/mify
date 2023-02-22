package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	removeClientName string
)

var removeClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Remove client",
	Long:  `Remove client from service`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.RemoveClient(appContext, workspacePath, ival, removeClientName); err != nil {
				fmt.Fprintf(os.Stderr, "failed to remove client to service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

var removeDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Remove database",
	Long:  `Remove database from service`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.RemovePostgres(appContext, workspacePath, ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to remove database from service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove client or database",
	Long:  `Remove client or database from service`,
	PersistentPreRun: func(*cobra.Command, []string) {
		err := appContext.LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init workspace: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
	removeClientCmd.PersistentFlags().StringVarP(&removeClientName, "to", "t", "", "Name of client service")
	if err := removeClientCmd.MarkPersistentFlagRequired("to"); err != nil {
		panic(err)
	}

	removeCmd.AddCommand(removeClientCmd)
	removeCmd.AddCommand(removeDatabaseCmd)
}
