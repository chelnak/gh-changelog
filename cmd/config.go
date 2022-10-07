// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"errors"
	"os"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/spf13/cobra"
)

var (
	output  string
	noColor bool
)

// configCmd is the entry point for printing the applications configuration in the terminal
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Prints the current configuration to the terminal in either json or yaml format. Defaults to yaml.",
	Long:  "Prints the current configuration to the terminal in either json or yaml format. Defaults to yaml.",
	RunE: func(command *cobra.Command, args []string) error {
		switch output {
		case "json":
			return configuration.Config.PrintJSON(noColor, os.Stdout)
		case "yaml":
			return configuration.Config.PrintYAML(noColor, os.Stdout)
		default:
			return errors.New("invalid output format. Valid values are 'json' and 'yaml'")
		}
	},
}

func init() {
	configCmd.Flags().StringVarP(&output, "output", "o", "yaml", "The output format. Valid values are 'json' and 'yaml'. Defaults to 'yaml'.")
	configCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable color output")
}
