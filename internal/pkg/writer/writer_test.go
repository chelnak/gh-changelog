package writer

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/stretchr/testify/assert"
)

type section struct {
	Name    string
	Message string
}

func TestWriter(t *testing.T) {

	tag := "v6.1.0"
	s := section{
		Name:    "Added",
		Message: "Added a new feature",
	}

	tests := []struct {
		name     string
		tag      string
		sections []section
		want     string
	}{
		{
			name: "starts with # Changelog",
		},
		{
			name:     "has the correct header tag",
			tag:      tag,
			sections: []section{s},
			want:     fmt.Sprintf("## \\[%s\\]", tag),
		},
		{
			name:     "contains a section header and entry",
			tag:      tag,
			sections: []section{s},
			want:     "### Added\n\n- Added a new feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			properties := changelog.NewChangeLogProperties("owner", "repo")

			tags := changelog.NewTagProperties(tt.tag, time.Now())

			for _, section := range tt.sections {
				tags.Append(section.Name, section.Message)
			}

			properties.Tags = append(properties.Tags, *tags)

			var output bytes.Buffer
			err := Write(properties, &output)
			assert.NoError(t, err)
			assert.Regexp(t, regexp.MustCompile(tt.want), output.String())
		})
	}

}
