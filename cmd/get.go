// Package cmd holds all top-level cobra commands. Each file should contain
// only one command and that command should have only one purpose.
package cmd

import (
	"fmt"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/get"
	"github.com/spf13/cobra"
)

var printLatest bool
var printVersion string

// getCmd retrieves a local changelog and prints it to stdout
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Parse a local changelog and print the result to stdout",
	Long:  "Parse a local changelog and print the result to stdout",
	RunE: func(command *cobra.Command, args []string) error {
		fileName := configuration.Config.FileName

		var changelog string
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

		fmt.Println(changelog)

		return nil
	},
}

func init() {

	getCmd.Flags().BoolVar(
		&printLatest,
		"latest",
		false,
		"Print the latest version of the changelog.",
	)

	getCmd.Flags().StringVar(
		&printVersion,
		"version",
		"",
		"A specific version to print.",
	)

	getCmd.Flags().SortFlags = false
}
