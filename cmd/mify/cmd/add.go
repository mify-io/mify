package cmd

import (
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	addClientName   string
	addFrontendType string
)

var addServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Add service",
	Long:  `Add service`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.CreateService(appContext, workspacePath, "go", ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

var addClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Add client",
	Long:  `Add client`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.AddClient(appContext, workspacePath, ival, addClientName); err != nil {
				fmt.Fprintf(os.Stderr, "failed to add client to service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

var addFrontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Add frontend",
	Long:  `Add frontend`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.CreateFrontend(appContext, workspacePath, "vue_js", ival); err != nil {
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
	addClientCmd.PersistentFlags().StringVarP(&addClientName, "to", "t", "", "Name of client service")
	addClientCmd.MarkPersistentFlagRequired("to")

	// TODO: limit witn enum
	addFrontendCmd.PersistentFlags().StringVarP(
		&addFrontendType,
		"template",
		"t",
		"vue",
		"Template (f.e. vue app)",
	)

	addCmd.AddCommand(addServiceCmd)
	addCmd.AddCommand(addClientCmd)
	addCmd.AddCommand(addFrontendCmd)
}
