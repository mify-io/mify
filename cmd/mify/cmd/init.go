package cmd

import (
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new workspace",
	Long:  `Initialize new workspace`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CreateWorkspace(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create workspace: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
}
