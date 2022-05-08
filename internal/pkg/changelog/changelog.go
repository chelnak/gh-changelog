package changelog

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	CurrentTag  string
	PreviousTag string
	Date        time.Time
	Added       []string
	Changed     []string
	Deprecated  []string
	Removed     []string
	Fixed       []string
	Security    []string
	Other       []string
}

func (e *Entry) append(section string, entry string) error {
	switch strings.ToLower(section) {
	case "added":
		e.Added = append(e.Added, entry)
	case "changed":
		e.Changed = append(e.Changed, entry)
	case "deprecated":
		e.Deprecated = append(e.Deprecated, entry)
	case "removed":
		e.Removed = append(e.Removed, entry)
	case "fixed":
		e.Fixed = append(e.Fixed, entry)
	case "security":
		e.Security = append(e.Security, entry)
	case "other":
		e.Other = append(e.Other, entry)
	default:
		return fmt.Errorf("unknown entry type '%s'", section)
	}

	return nil
}

type Changelog interface {
	GetRepoName() string
	GetRepoOwner() string
	GetUnreleased() []string
	GetEntries() []Entry
}

type changelog struct {
	repoName   string
	repoOwner  string
	unreleased []string
	entries    []Entry
}

func (c *changelog) GetRepoName() string {
	return c.repoName
}

func (c *changelog) GetRepoOwner() string {
	return c.repoOwner
}

func (c *changelog) GetEntries() []Entry {
	return c.entries
}

func (c *changelog) GetUnreleased() []string {
	return c.unreleased
}
