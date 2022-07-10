//Package writer is responsible for parsing the given changelog struct
//into a go template and writing it to the given writer.
package writer

import (
	"io"
	"text/template"

	"github.com/chelnak/gh-changelog/internal/changelog"
)

var tmplSrc = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

{{- $unreleased := .GetUnreleased }}
{{- if $unreleased }}

## Unreleased
{{range $unreleased }}
- {{.}}
{{- end}}
{{- end -}}
{{range .GetEntries}}
## [{{.CurrentTag}}](https://github.com/{{$.GetRepoOwner}}/{{$.GetRepoName}}/tree/{{.CurrentTag}}) - {{.Date.Format "2006-01-02"}}

[Full Changelog](https://github.com/{{$.GetRepoOwner}}/{{$.GetRepoName}}/compare/{{.PreviousTag}}...{{.CurrentTag}})
{{if .Added }}
### Added
{{range .Added}}
- {{.}}
{{- end}}
{{- end -}}
{{- if .Changed }}
### Changed
{{- range .Changed}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Deprecated }}
### Deprecated
{{range .Deprecated}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Removed }}
### Removed
{{range .Removed}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Fixed }}
### Fixed
{{range .Fixed}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Security }}
### Security
{{range .Security}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Other }}
### Other
{{range .Other}}
- {{.}}
{{- end}}
{{- end}}
{{- end}}`

func Write(writer io.Writer, changelog changelog.Changelog) error {
	tmpl := template.Must(template.New("changelog").Parse(tmplSrc))
	return tmpl.Execute(writer, changelog)
}
