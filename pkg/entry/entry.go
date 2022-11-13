// Package entry provides a datastructure for changelog entries.
package entry

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Entry represents a single entry in the changelog
type Entry struct {
	Previous *Entry // Get or Set the previous entry in the changelog.
	Next     *Entry // Get or Set the next entry in the changelog.

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

// Append updates the given section in the entry..
func (e *Entry) Append(section string, entry string) error {
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

// GetSection uses reflection to return a given section in the entry.
// If the section does not exist, an empty slice is returned.
func (e *Entry) GetSection(section string) []string {
	title := cases.Title(language.English)
	ref := reflect.ValueOf(e).Elem().FieldByName(title.String(section))
	if ref.IsValid() {
		return ref.Interface().([]string)
	}
	return nil
}

// NewEntry creates a new entry (node) that can be added to the changelog datastructure.
func NewEntry(tag string, date time.Time) Entry {
	return Entry{
		Tag:  tag,
		Date: date,
	}
}
