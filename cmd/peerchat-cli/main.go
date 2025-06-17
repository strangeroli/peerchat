package main

import (
	"os"

	"github.com/Xelvra/peerchat/internal/cli"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	version = "0.2.0-alpha"
)

func main() {
	Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd := cli.CreateRootCommand(version)

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xelvra/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Initialize configuration
	cobra.OnInitialize(initConfig)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Configuration loading temporarily disabled for debugging
	// TODO: Implement configuration loading
}
