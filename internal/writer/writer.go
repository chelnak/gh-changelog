// Package writer is responsible for parsing the given changelog struct
// into a go template and writing it to the given writer.
package writer

import (
	"io"
	"os/exec"
	"text/template"

	"github.com/chelnak/gh-changelog/internal/gitclient"
	"github.com/chelnak/gh-changelog/pkg/changelog"
)

var tmplSrc = `<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

{{- $unreleased := .GetUnreleased }}
{{if $unreleased }}
## Unreleased
{{range $unreleased }}
- {{.}}
{{- end}}
{{- end -}}
{{range .GetEntries}}
## [{{.Tag}}](https://github.com/{{$.GetRepoOwner}}/{{$.GetRepoName}}/tree/{{.Tag}}) - {{.Date.Format "2006-01-02"}}
{{ if .Previous }}
[Full Changelog](https://github.com/{{$.GetRepoOwner}}/{{$.GetRepoName}}/compare/{{.Previous.Tag}}...{{.Tag}})
{{else}}
[Full Changelog](https://github.com/{{$.GetRepoOwner}}/{{$.GetRepoName}}/compare/{{if .PrevTag }}{{.PrevTag}}{{else}}{{getFirstCommit}}{{end}}...{{.Tag}})
{{- end -}}

{{- if .Security }}
### Security
{{range .Security}}
- {{.}}
{{- end}}
{{end}}
{{- if .Changed }}
### Changed
{{range .Changed}}
- {{.}}
{{- end}}
{{end}}
{{- if .Removed }}
### Removed
{{range .Removed}}
- {{.}}
{{- end}}
{{end}}
{{- if .Deprecated }}
### Deprecated
{{range .Deprecated}}
- {{.}}
{{- end}}
{{end}}
{{- if .Added }}
### Added
{{range .Added}}
- {{.}}
{{- end}}
{{end}}
{{- if .Fixed }}
### Fixed
{{range .Fixed}}
- {{.}}
{{- end}}
{{end}}
{{- if .Other }}
### Other
{{range .Other}}
- {{.}}
{{- end}}
{{end}}
{{- end}}
`

func Write(writer io.Writer, changelog changelog.Changelog) error {
	tmpl, err := template.New("changelog").Funcs(template.FuncMap{
		"getFirstCommit": func() string {
			git := gitclient.NewGitClient(exec.Command)
			commit, err := git.GetFirstCommit()
			if err != nil {
				return ""
			}
			return commit
		},
	}).Parse(tmplSrc)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, changelog)
}
