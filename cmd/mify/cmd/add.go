package cmd

import (
	"fmt"
	"os"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"
)

var (
	addClientName      string
	addServiceLanguage string
	addFrontendType    string
	addDatabaseEngine  string
)

var addServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Add service",
	Long:  `Add a new service in workspace`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.CreateService(appContext, workspacePath, addServiceLanguage, ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

var addClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Add client",
	Long:  `Add a client from one service or frontend to another`,
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
	Long:  `Add a frontend to workspace`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.CreateFrontend(appContext, workspacePath, addFrontendType, ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

var addApiGatewayCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "Add api gateway",
	Long:  `Add an api gateway to workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mify.CreateApiGateway(appContext); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create api gateway: %s\n", err)
			os.Exit(2)
		}
	},
}

var addDatabase = &cobra.Command{
	Use:   "database",
	Short: "Add database",
	Long:  `Add a database to workspace (default is Postgres)`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, ival := range args {
			if err := mify.AddPostgres(appContext, workspacePath, ival); err != nil {
				fmt.Fprintf(os.Stderr, "failed to add postgres to service: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add <service|client|frontend|database>",
	Long:  `Add a service, frontend or clients, or database`,
	PersistentPreRun: func(*cobra.Command, []string) {
		err := appContext.LoadWorkspace()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init workspace: %s\n", err)
			os.Exit(2)
		}
	},
}

func init() {
	addClientCmd.PersistentFlags().StringVarP(&addClientName, "to", "t", "", "Name of client service")
	if err := addClientCmd.MarkPersistentFlagRequired("to"); err != nil {
		panic(err)
	}
	if err := addClientCmd.MarkPersistentFlagRequired("to"); err != nil {
		panic(err)
	}
	addServiceCmd.PersistentFlags().StringVarP(
		&addServiceLanguage,
		"language", "l", "go", "Choose language for service: go, python",
	)

	// TODO: limit witn enum
	addFrontendCmd.PersistentFlags().StringVarP(
		&addFrontendType,
		"template",
		"t",
		"nuxtjs",
		"Template (e.g. nuxtjs, react-ts)",
	)

	addDatabase.PersistentFlags().StringVarP(
		&addDatabaseEngine,
		"engine",
		"e",
		"postgres",
		"DB Engine (Postgres, for now)",
	)

	addCmd.AddCommand(addServiceCmd)
	addCmd.AddCommand(addClientCmd)
	addCmd.AddCommand(addFrontendCmd)
	addCmd.AddCommand(addApiGatewayCmd)
	addCmd.AddCommand(addApiGatewayCmd)
	addCmd.AddCommand(addDatabase)
}
