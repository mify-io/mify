package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/mify-io/mify/internal/mify"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile       string
	workspacePath string
	isVerbose     bool
	appContext    *mify.CliContext
)

var MIFY_VERSION string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mify",
	Short: "mify CLI tool",
	Long: `Code generation of services across your repository.
	The available commands for execution are listed below.`,
	Version: MIFY_VERSION,
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
	cleanup()
}

func cleanup() {
	if appContext == nil {
		return
	}
	appContext.Cancel()

	if err := mify.Cleanup(appContext); err != nil {
		fmt.Fprintf(os.Stderr, "failed to cleanup: %s", err)
		os.Exit(2)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&workspacePath, "path", "p", "", "Path to workspace")
	rootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Show verbose output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(toolCmd)
	rootCmd.AddCommand(cloudCmd)

	rootCmd.PersistentPostRun = PersistentPostRun
}

func PersistentPostRun(cmd *cobra.Command, args []string) {
	if strings.HasPrefix(cmd.Name(), "__") {
		// ignore __complete and other
		return
	}

	desc := appContext.GetWorkspaceDescription()
	if desc == nil {
		// TODO: we should handle commands without workspace (such as completion generation). Move queue file to linux tmp file?
		return
	}

	appContext.InitStatsCollector(desc.GetStatsQueueFile())
	appContext.StatsCollector.LogCobraCommandExecuted(cmd)
	err := appContext.StatsCollector.MaybeSendStats()
	if err != nil {
		appContext.Logger.Printf("can't send stats: %s", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(mify.GetConfigDirectory())
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	config := mify.NewDefaultConfig()
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config: %s", err)
		os.Exit(2)
	}
	appContext = mify.NewContext(config, workspacePath, isVerbose, MIFY_VERSION)
}
