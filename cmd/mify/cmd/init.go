package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	vcsTemplate string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new workspace",
	Long:  `Initialize new workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		appContext.StatsCollector.LogEvent("run", cmd)
		workspaceName := "."
		if len(args) > 0 {
			workspaceName = args[0]
		}

		if err := mify.CreateWorkspace(appContext, workspacePath, workspaceName, vcsTemplate); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create workspace: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
	initCmd.LocalFlags().StringVar(&vcsTemplate, "vcs", "git", "[git|none]")
}
