package cmd

import (
	"os"

	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/chelnak/gh-changelog/internal/pkg/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nextVersion string

// newCmd is the entry point for creating a new changelog
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new changelog from activity in the current repository",
	Long:  "Creates a new changelog from activity in the current repository.",
	RunE: func(command *cobra.Command, args []string) error {
		builder := changelog.NewChangelogBuilder()
		builder = builder.WithSpinner(true)
		builder = builder.WithNextVersion(nextVersion)

		changelog, err := builder.Build()
		if err != nil {
			return err
		}

		fileName := viper.GetString("file_name")
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		return writer.Write(f, changelog)
	},
}

func init() {
	newCmd.Flags().StringVar(&nextVersion, "next-version", "", "The next version to be released. The value passed does not have to be an existing tag.")
}
