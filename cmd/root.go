//Package cmd holds all top-level cobra commands. Each file should contain
//only one command and that command should have only one purpose.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/chelnak/gh-changelog/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var version = "dev"
var errSilent = errors.New("ErrSilent")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "changelog [command]",
	Short:         "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Long:          "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run:           nil,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if configuration.Config.CheckForUpdates {
			utils.CheckForUpdate(version)
		}
	},
}

func init() {
	err := configuration.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errSilent
	})

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(configCmd)
}

// Execute is called from main and is responsible for processing
// requests to the application and handling exit codes appropriately
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if err != errSilent {
			fmt.Fprintln(os.Stderr, fmt.Errorf("❌ %s", err))
		}
		return 1
	}
	return 0
}
