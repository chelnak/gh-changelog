// package get retrieves a local changelog, parses it and returns the result.
package get

import (
	"bytes"

	"github.com/chelnak/gh-changelog/internal/writer"
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

func writeBuf(parsedChangelog changelog.Changelog) (string, error) {
	var buf bytes.Buffer
	if err := writer.Write(&buf, parsedChangelog); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func changelogWithSingleEntry(entry entry.Entry, repoName, repoOwner string) changelog.Changelog {
	// Isolate the entry
	entry.Next = nil
	entry.Previous = nil

	cl := changelog.NewChangelog(repoName, repoOwner)
	cl.Insert(entry)
	return cl
}

// GetVersion retrieves a local changelog, parses it and returns a string
// containing only the specified version.
func GetVersion(fileName string, tag string) (string, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return "", err
	}

	versionEntry := getTag(parsedChangelog, tag)
	cl := changelogWithSingleEntry(
		*versionEntry,
		parsedChangelog.GetRepoName(),
		parsedChangelog.GetRepoOwner(),
	)

	return writeBuf(cl)
}

// GetLatest retrieves a local changelog, parses it and returns a string
// containing only the latest entry.
func GetLatest(fileName string) (string, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return "", err
	}

	latestEntry := parsedChangelog.Tail()
	cl := changelogWithSingleEntry(
		*latestEntry,
		parsedChangelog.GetRepoName(),
		parsedChangelog.GetRepoOwner(),
	)

	return writeBuf(cl)
}

// Get retrieves a local changelog, parses it and returns a string
// containing all entries
func GetAll(fileName string) (string, error) {
	parsedChangelog, err := parseChangelog(fileName)
	if err != nil {
		return "", err
	}

	return writeBuf(parsedChangelog)
}
