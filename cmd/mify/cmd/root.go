package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile       string
	workspacePath string
	appContext    *mify.CliContext
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mify",
	Short: "mify CLI tool",
	Long: `Code generation of services across your repository.
	The available commands for execution are listed below.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var endWaiter sync.WaitGroup
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		endWaiter.Add(1)
		cleanup()
		endWaiter.Done()
	}()

	cobra.CheckErr(rootCmd.Execute())
	endWaiter.Wait()
}

func cleanup() {
	appContext.Cancel()

	if err := mify.Cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to cleanup: %s", err)
		os.Exit(2)
	}
}

func init() {
	appContext = mify.NewContext()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&workspacePath, "path", "p", "", "Path to workspace")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(runCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mify" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
