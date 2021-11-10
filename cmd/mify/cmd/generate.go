package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate",
	Long:  `generate`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.ServiceGenerate(appContext, workspacePath, ival); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Fprintf(os.Stderr, "failed to generate in service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

func init() {
}
