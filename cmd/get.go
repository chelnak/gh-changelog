// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"bytes"
	"fmt"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/get"
	"github.com/chelnak/gh-changelog/internal/writer"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/spf13/cobra"
)

var printLatest bool
var printVersion string

// getCmd retrieves a local changelog and prints it to stdout
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Reads a changelog file and prints the result to stdout",
	Long: `Reads a changelog file and prints the result to stdout.

This command is useful for creating and updating Release notes in GitHub.

┌─────────────────────────────────────────────────────────────────────┐
│Example                                                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│→ gh changelog get --latest > release_notes.md                       │
│→ gh release create --title "Release v1.0.0" -F release_notes.md     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
`,
	RunE: func(command *cobra.Command, args []string) error {
		fileName := configuration.Config.FileName

		var changelog changelog.Changelog
		var err error

		if printLatest {
			changelog, err = get.GetLatest(fileName)
		} else if printVersion != "" {
			changelog, err = get.GetVersion(fileName, printVersion)
		} else {
			changelog, err = get.GetAll(fileName)
		}

		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := writer.Write(&buf, changelog); err != nil {
			return err
		}

		fmt.Println(buf.String())

		return nil
	},
}

func init() {
	getCmd.Flags().BoolVar(
		&printLatest,
		"latest",
		false,
		"Prints the latest version from the changelog to stdout.",
	)

	getCmd.Flags().StringVar(
		&printVersion,
		"version",
		"",
		"Prints a specific version from the changelog to stdout.",
	)

	getCmd.Flags().SortFlags = false
}
