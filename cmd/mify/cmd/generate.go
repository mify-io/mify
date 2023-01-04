package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	migrate bool
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "generate [service]...",
	Short: "Generate code in workspace",
	Long:  `Generate code for given list of services after schema changes.`,
	PersistentPreRun: func(*cobra.Command, []string) {
		err := appContext.LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init workspace: %s\n", err)
			os.Exit(2)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		appContext.StatsCollector.LogEvent("run", cmd)
		if err := mify.ServiceGenerateMany(appContext, workspacePath, args, migrate); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(2)
		}
	},
}

func init() {
	genCmd.Flags().BoolVarP(
		&migrate,
		"migrate", "m", true, "Should code migrations be applied?",
	)
}
