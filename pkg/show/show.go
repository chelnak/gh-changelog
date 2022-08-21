// Package show is responsible for rendering the contents of a
// given CHANGELOG.md file and displaying it in the terminal.
package show

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/chelnak/gh-changelog/internal/viewport"
)

func render(content string) error {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(140))

	content, err := r.Render(content)
	if err != nil {
		return err
	}

	return viewport.Start(content)
}

func RenderFile(path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return errors.New("changelog not found. Check your configuration or run gh changelog new")
	}

	return render(string(data))
}

func RenderString(content string) error {
	if content == "" {
		return errors.New("there is no content to render")
	}

	return render(content)
}
