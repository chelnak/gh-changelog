package markdown

import (
	"errors"
	"io/ioutil"

	"github.com/charmbracelet/glamour"
)

func Render(path string) error {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(140), // TODO: make this configurable
	)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("❌ changelog not found. Check your configuration or run gh changelog new")
	}

	content, err := r.Render(string(data))
	if err != nil {
		return err
	}

	return start(content)
}
