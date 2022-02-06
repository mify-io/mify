package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var deployCmd = &cobra.Command{
	Use:   "deploy [service]...",
	Short: "Run deploy tool",
	Long:  `Run deploy tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.Deploy(appContext, workspacePath, args); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run deploy: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
}
