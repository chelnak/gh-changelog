package writer

import (
	"os"
	"text/template"

	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/spf13/viper"
)

//lintLint:ignore U1000
func Write(changeLog *changelog.ChangeLogProperties) error {
	var tmplSrc = `# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).
{{range .Tags}}
## [{{.Tag}}](https://github.com/{{$.RepoOwner}}/{{$.RepoName}}/tree/{{.Tag}}) - ({{.Date.Format "2006-01-02"}})

[Full Changelog](https://github.com/{{$.RepoOwner}}/{{$.RepoName}}/compare/v6.1.0...v6.1.1)


{{- if .Added }}
### Added
{{range .Added}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Changed }}
### Changed
{{range .Changed}}
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

	tmpl := template.Must(template.New("changelog").Parse(tmplSrc))

	fileName := viper.GetString("fileName")
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, changeLog)
	if err != nil {
		return err
	}
	return nil
}
