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
	Long:  `Remove client`,
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
	Short: "remove",
	Long:  `remove`,
}

func init() {
	removeClientCmd.PersistentFlags().StringVarP(&removeClientName, "to", "t", "", "Name of client service")
	removeClientCmd.MarkPersistentFlagRequired("to")

	removeCmd.AddCommand(removeClientCmd)
}
