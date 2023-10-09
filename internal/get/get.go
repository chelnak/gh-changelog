// Package get retrieves a local changelog, parses it and returns the result.
package get

import (
	"fmt"

	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/chelnak/gh-changelog/pkg/parser"
)

func getTag(parsedChangelog changelog.Changelog, tag string) *entry.Entry {
	currentEntry := parsedChangelog.Head()
	for currentEntry != nil {
		if currentEntry.Tag == tag {
			return currentEntry
		}
		currentEntry = currentEntry.Next
	}
	return nil
}

func parseChangelog(fileName string) (changelog.Changelog, error) {
	parser := parser.NewParser(fileName, "", "")
	return parser.Parse()
}

func changelogWithSingleEntry(entry entry.Entry, repoName, repoOwner string) changelog.Changelog {
	// Isolate the entry
	entry.Next = nil
	entry.Previous = nil

	cl := changelog.NewChangelog(repoOwner, repoName)
	cl.Insert(entry)
	return cl
}

// GetVersion retrieves a local changelog, parses it and returns a string
// containing only the specified version.
func GetVersion(fileName string, tag string) (changelog.Changelog, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return nil, err
	}

	versionEntry := getTag(parsedChangelog, tag)
	if versionEntry == nil {
		return nil, fmt.Errorf("version %s not found", tag)
	}

	cl := changelogWithSingleEntry(
		*versionEntry,
		parsedChangelog.GetRepoName(),
		parsedChangelog.GetRepoOwner(),
	)

	return cl, nil
}

// GetLatest retrieves a local changelog, parses it and returns a string
// containing only the latest entry.
func GetLatest(fileName string) (changelog.Changelog, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return nil, err
	}

	latestEntry := parsedChangelog.Tail()
	cl := changelogWithSingleEntry(
		*latestEntry,
		parsedChangelog.GetRepoName(),
		parsedChangelog.GetRepoOwner(),
	)

	return cl, nil
}

// GetAll retrieves a local changelog, parses it and returns a string
// containing all entries
func GetAll(fileName string) (changelog.Changelog, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return nil, err
	}

	return parsedChangelog, nil
}
