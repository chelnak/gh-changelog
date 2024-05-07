// Package changelog provides the data structure that is responsible for
// holding the changelog data.
package changelog

import (
	"github.com/chelnak/gh-changelog/pkg/entry"
)

type Changelog struct {
	// Support for linked list structure
	head *entry.Entry
	tail *entry.Entry

	repoName   string
	repoOwner  string
	unreleased []string
}

// NewChangelog creates a new changelog data structure.
func NewChangelog(repoOwner string, repoName string) *Changelog {
	return &Changelog{
		repoName:   repoName,
		repoOwner:  repoOwner,
		unreleased: []string{},
	}
}

// GetRepoName returns the name of the repository.
func (c *Changelog) GetRepoName() string {
	return c.repoName
}

// GetRepoOwner returns the owner of the repository.
func (c *Changelog) GetRepoOwner() string {
	return c.repoOwner
}

// GetUnreleased returns the unreleased changes if any exist.
func (c *Changelog) GetUnreleased() []string {
	return c.unreleased
}

// AddUnreleased adds a list of unreleased changes to the changelog.
// This only needs to be a slice of strings.
func (c *Changelog) AddUnreleased(entry []string) {
	c.unreleased = append(c.unreleased, entry...)
}

// Insert inserts a new entry into the changelog.
func (c *Changelog) Insert(e entry.Entry) {
	if c.head != nil {
		e.Next = c.head
		c.head.Previous = &e
	}
	c.head = &e

	currentEntry := c.head
	for currentEntry.Next != nil {
		currentEntry = currentEntry.Next
	}
	c.tail = currentEntry
}

// GetEntries returns a list of entries in the changelog.
// This is a convenience method that creates a contiguous list of entries
// that can be iterated over. The latest entry will be the first item in the list.
func (c *Changelog) GetEntries() []*entry.Entry {
	var entries []*entry.Entry
	currentEntry := c.tail
	for currentEntry != nil {
		entries = append(entries, currentEntry)
		currentEntry = currentEntry.Previous
	}

	return entries
}

// Head returns the first entry in the changelog.
func (c *Changelog) Head() *entry.Entry {
	return c.head
}

// Tail returns the last entry in the changelog.
func (c *Changelog) Tail() *entry.Entry {
	return c.tail
}
