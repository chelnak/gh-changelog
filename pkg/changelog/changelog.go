// Package changelog provides the datastructure that is responsible for
// holding the changelog data.
package changelog

import (
	"github.com/chelnak/gh-changelog/pkg/entry"
)

// Changelog is an interface for a changelog datastructure.
type Changelog interface {
	GetRepoName() string
	GetRepoOwner() string
	GetUnreleased() []string
	AddUnreleased([]string)
	Insert(entry.Entry)
	GetEntries() []*entry.Entry
	Head() *entry.Entry
	Tail() *entry.Entry
}

type changelog struct {
	// Support for linked list structure
	head *entry.Entry
	tail *entry.Entry

	repoName   string
	repoOwner  string
	unreleased []string
}

// GetRepoName returns the name of the repository.
func (c *changelog) GetRepoName() string {
	return c.repoName
}

// GetRepoOwner returns the owner of the repository.
func (c *changelog) GetRepoOwner() string {
	return c.repoOwner
}

// GetUnreleased returns the unreleased changes if any exist.
func (c *changelog) GetUnreleased() []string {
	return c.unreleased
}

// AddUnreleased adds a list of unreleased changes to the changelog.
// This only needs to be a slice of strings.
func (c *changelog) AddUnreleased(entry []string) {
	c.unreleased = append(c.unreleased, entry...)
}

// Insert inserts a new entry into the changelog.
func (c *changelog) Insert(e entry.Entry) {
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
// that can be iterated over. The latest entry will be the first item in the list..
func (c *changelog) GetEntries() []*entry.Entry {
	entries := []*entry.Entry{}
	currentEntry := c.tail
	for currentEntry != nil {
		entries = append(entries, currentEntry)
		currentEntry = currentEntry.Previous
	}

	return entries
}

// Head returns the first entry in the changelog.
func (c *changelog) Head() *entry.Entry {
	return c.head
}

// Tail returns the last entry in the changelog.
func (c *changelog) Tail() *entry.Entry {
	return c.tail
}

// NewChangelog creates a new changelog datastructure.
func NewChangelog(repoOwner string, repoName string) Changelog {
	return &changelog{
		repoName:   repoName,
		repoOwner:  repoOwner,
		unreleased: []string{},
	}
}
