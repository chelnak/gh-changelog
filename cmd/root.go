package cmd

import (
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

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "changelog [args]",
	Short:   "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Long:    "Create a changelog that adheres to the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) format",
	Version: version,
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
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
