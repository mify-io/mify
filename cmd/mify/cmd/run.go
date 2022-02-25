package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all go services",
	Long:  `Run all go services in one dev-runner process`,
	PersistentPreRun: func(*cobra.Command, []string) {
		err := appContext.LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init workspace: %s\n", err)
			os.Exit(2)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.Run(appContext, workspacePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run app: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
}
