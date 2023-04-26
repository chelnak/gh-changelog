// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/utils"
	"github.com/spf13/cobra"
)

var version = "dev"
var errSilent = errors.New("ErrSilent")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "changelog [command]",
	Short: "A GitHub CLI extension that will make your changelogs ✨",
	Long: `gh changelog - A GitHub CLI extension that will make your changelogs ✨

Easily create standardised changelogs for your project that follow
conventions set by the keepachangelog project.

For more information check out the following link:

🔗 https://keepachangelog.com

Getting started is easy:

┌────────────────────┐
│•••                 │
├────────────────────┤
│                    │
│→ gh changelog new  │
└────────────────────┘

You can also view the changelog at any time:

┌────────────────────┐
│•••                 │
├────────────────────┤
│                    │
│→ gh changelog show │
└────────────────────┘

Issues or feature requests can be opened at:

🔗 https://github.com/chelnak/gh-changelog/issues`,
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
		formatError(err)
		os.Exit(1)
	}

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errSilent
	})

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(parseCmd)
}

func formatError(err error) {
	fmt.Print("\n❌ It looks like something went wrong!\n")
	fmt.Println("\nReported errors:")
	fmt.Fprintln(os.Stderr, fmt.Errorf("• %s", err))
	fmt.Println()
}

// Execute is called from main and is responsible for processing
// requests to the application and handling exit codes appropriately
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if err != errSilent {
			formatError(err)
		}
		return 1
	}
	return 0
}
