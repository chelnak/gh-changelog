package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/chelnak/gh-changelog/internal/pkg/writer"
	"github.com/spf13/cobra"
)

var version = "dev"
var ErrSilent = errors.New("ErrSilent")

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "changelog",
	Short:         "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Long:          "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(command *cobra.Command, args []string) error {

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		_ = s.Color("green")
		s.FinalMSG = "✅ Done!\n"

		changeLog, err := changelog.MakeFullChangelog(s)
		if err != nil {
			return err
		}

		s.Stop()

		return writer.Write(changeLog)
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
		return ErrSilent
	})
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
