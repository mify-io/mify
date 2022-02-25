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
				fmt.Fprintf(os.Stderr, "failed to add client to service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

// addCmd represents the add command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove client",
	Long:  `Remove client from service`,
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
}
