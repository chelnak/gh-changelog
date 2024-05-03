// Package parser provides a simple interface for parsing markdown changelogs.
package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chelnak/gh-changelog/internal/utils"
	"github.com/chelnak/gh-changelog/pkg/changelog"
	"github.com/chelnak/gh-changelog/pkg/entry"
	"github.com/gomarkdown/markdown/ast"
	mdparser "github.com/gomarkdown/markdown/parser"
)

type parser struct {
	path      string
	repoOwner string
	repoName  string
}

// Parser is an interface for parsing markdown changelogs.
type Parser interface {
	Parse() (changelog.Changelog, error)
}

// NewParser returns a new parser for the given changelog..
func NewParser(path, repoOwner, repoName string) Parser {
	return &parser{
		path:      path,
		repoName:  repoName,
		repoOwner: repoOwner,
	}
}

// Parse parses the changelog and returns a Changelog struct.
func (p *parser) Parse() (changelog.Changelog, error) {
	repoContext, err := utils.GetRepoContext()
	if err != nil {
		return nil, err
	}

	if p.repoOwner == "" {
		p.repoOwner = repoContext.Owner
	}

	if p.repoName == "" {
		p.repoName = repoContext.Name
	}

	data, err := os.ReadFile(filepath.Clean(p.path))
	if err != nil {
		return nil, err
	}

	markdownParser := mdparser.New()
	output := markdownParser.Parse(data)

	var tagIndex []string // This is a list of tags in order
	var unreleased []string
	var entries = map[string]*entry.Entry{} // Maintain a map of tag to entry
	var currentTag string
	var currentSection string

	for _, child := range output.GetChildren() {
		switch child.(type) {
		case *ast.Heading:
			if isHeading(child, 2) {
				currentTag = getTagFromHeadingLink(child)
				if currentTag == "" && isHeadingUnreleased(child) {
					currentTag = "Unreleased"
					continue
				}
				date := getDateFromHeading(child)
				if _, ok := entries[currentTag]; !ok {
					e := entry.NewEntry(currentTag, date)
					entries[currentTag] = &e
					tagIndex = append(tagIndex, currentTag)
				}
			}

			if isHeading(child, 3) {
				currentSection = getTextFromChildNodes(child)
			}
		case *ast.List:
			items := getItemsFromList(child)
			if currentTag == "Unreleased" {
				for _, item := range items {
					unreleased = append(unreleased, getTextFromChildNodes(item))
				}
				continue
			}

			for _, item := range items {
				err := entries[currentTag].Append(currentSection, getTextFromChildNodes(item))
				if err != nil {
					// TODO: Add more context to this error
					return nil, fmt.Errorf("error parsing changelog: %s", err)
				}
			}
		default:
			// TODO: Add more context to this block
			// We are ignoring other types of nodes for now
			continue
		}
	}

	cl := changelog.NewChangelog(p.repoOwner, p.repoName)

	if len(unreleased) > 0 {
		cl.AddUnreleased(unreleased)
	}

	for _, tag := range tagIndex {
		cl.Insert(*entries[tag])
	}

	return cl, nil
}

func isListItem(node ast.Node) bool {
	_, ok := node.(*ast.ListItem)
	return ok
}

func isLink(node ast.Node) bool {
	_, ok := node.(*ast.Link)
	return ok
}

func isText(node ast.Node) bool {
	_, ok := node.(*ast.Text)
	return ok
}

func isHeading(node ast.Node, level int) bool {
	_, heading := node.(*ast.Heading)
	ok := heading && node.(*ast.Heading).Level == level
	return ok
}

func isParagraph(node ast.Node) bool {
	_, ok := node.(*ast.Paragraph)
	return ok
}

func getItemsFromList(node ast.Node) []*ast.ListItem {
	var items []*ast.ListItem
	for _, child := range node.GetChildren() {
		if isListItem(child) {
			items = append(items, child.(*ast.ListItem))
		}
	}
	return items
}

func getTextFromChildNodes(node ast.Node) string {
	var text []string
	for _, child := range node.GetChildren() {
		if isParagraph(child) {
			text = append(text, getTextFromChildNodes(child)) // This stinks
			// text = append(text, "\n")                         // so does this
		}

		if isText(child) {
			text = append(text, string(child.(*ast.Text).Literal))
		}

		if isLink(child) {
			linkText := getTextFromChildNodes(child)
			link := fmt.Sprintf("[%s](%s)", linkText, child.(*ast.Link).Destination)
			text = append(text, link)
		}
	}
	return strings.Join(text, "")
}

func getTagFromHeadingLink(node ast.Node) string {
	for _, child := range node.GetChildren() {
		if isLink(child) {
			return getTextFromChildNodes(child)
		}
	}
	return ""
}

func getDateFromHeading(node ast.Node) time.Time {
	var date time.Time
	for _, child := range node.GetChildren() {
		if isText(child) {
			text := string(child.(*ast.Text).Literal)
			if text != "" {
				date, err := time.Parse("2006-01-02", strings.ReplaceAll(text, " - ", ""))
				if err != nil {
					panic(err)
				}
				return date
			}
		}
	}
	return date
}

func isHeadingUnreleased(node ast.Node) bool {
	return strings.Contains(getTextFromChildNodes(node), "Unreleased")
}
