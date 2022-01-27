package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "generate [service]",
	Short: "generate",
	Long:  `generate`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: make service name as optional
		for _, ival := range args {
			if err := mify.ServiceGenerate(appContext, workspacePath, ival); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Fprintf(os.Stderr, "service '%s' generation failed: %s", ival, err)
				os.Exit(2)
			}
		}
	},
}

func init() {
}
