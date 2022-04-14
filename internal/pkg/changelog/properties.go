package changelog

import (
	"fmt"
	"strings"
	"time"
)

type ChangeLogProperties struct {
	RepoName  string
	RepoOwner string
	Tags      []TagProperties
}

func NewChangeLogProperties(repoName string, repoOwner string) *ChangeLogProperties {
	return &ChangeLogProperties{
		RepoName:  repoName,
		RepoOwner: repoOwner,
		Tags:      []TagProperties{},
	}
}

type TagProperties struct {
	Tag        string
	Date       time.Time
	Added      []string
	Changed    []string
	Deprecated []string
	Removed    []string
	Fixed      []string
	Security   []string
	Other      []string
}

func NewTagProperties(tag string, date time.Time) *TagProperties {
	return &TagProperties{
		Tag:        tag,
		Date:       date,
		Added:      []string{},
		Changed:    []string{},
		Deprecated: []string{},
		Removed:    []string{},
		Fixed:      []string{},
		Security:   []string{},
		Other:      []string{},
	}
}

func (tp *TagProperties) Append(entryType string, entry string) error {
	switch strings.ToLower(entryType) {
	case "added":
		tp.Added = append(tp.Added, entry)
	case "changed":
		tp.Changed = append(tp.Changed, entry)
	case "deprecated":
		tp.Deprecated = append(tp.Deprecated, entry)
	case "removed":
		tp.Removed = append(tp.Removed, entry)
	case "fixed":
		tp.Fixed = append(tp.Fixed, entry)
	case "security":
		tp.Security = append(tp.Security, entry)
	case "other":
		tp.Other = append(tp.Other, entry)
	default:
		return fmt.Errorf("unknown entry type: %s", entryType)
	}

	return nil
}
