package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/spf13/cobra"
)

var version = "dev"
var ErrSilent = errors.New("ErrSilent")

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "changelog [command]",
	Short:         "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Long:          "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run:           nil,
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
		return ErrSilent
	})

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(showCmd)
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if err != ErrSilent {
			fmt.Fprintln(os.Stderr, err)
		}
		return 1
	}
	return 0
}
