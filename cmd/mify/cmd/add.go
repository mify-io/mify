package cmd

import (
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify"
	"github.com/spf13/cobra"
)

var addServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Add service",
	Long:  `Add service`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.CreateService(workspacePath, ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  `add`,
}

func init() {
	addCmd.AddCommand(addServiceCmd)
}
