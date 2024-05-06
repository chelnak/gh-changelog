package githubclient

import (
	"fmt"
	"sort"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type RefNode struct {
	Name   string
	Target struct {
		TypeName string `graphql:"__typename"`
		Tag      struct {
			Oid    string
			Tagger struct {
				Date time.Time
			}
		} `graphql:"... on Tag"`
		Commit struct {
			Oid       string
			Committer struct {
				Date time.Time
			}
		} `graphql:"... on Commit"`
	}
}

type TagQuery struct {
	Repository struct {
		Refs struct {
			Nodes    []RefNode
			PageInfo struct {
				EndCursor   graphql.String
				HasNextPage bool
			}
		} `graphql:"refs(refPrefix: \"refs/tags/\", last: 100, after: $cursor)"`
	} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
}

type Tag struct {
	Name string
	Sha  string
	Date time.Time
}

func (client *GitHub) GetTags() ([]Tag, error) {
	variables := map[string]interface{}{
		"repositoryOwner": graphql.String(client.repoContext.owner),
		"repositoryName":  graphql.String(client.repoContext.name),
		"cursor":          (*graphql.String)(nil),
	}

	var tags []Tag
	var tagQuery TagQuery
	var nodes []RefNode

	for {
		err := client.base.Query("tags", &tagQuery, variables)
		if err != nil {
			return nil, fmt.Errorf("error getting tags: %w", err)
		}

		nodes = append(nodes, tagQuery.Repository.Refs.Nodes...)

		if !tagQuery.Repository.Refs.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = tagQuery.Repository.Refs.PageInfo.EndCursor
	}

	for _, node := range nodes {
		switch node.Target.TypeName {
		case "Tag":
			tags = append(tags, Tag{
				Name: node.Name,
				Sha:  node.Target.Tag.Oid,
				Date: node.Target.Tag.Tagger.Date,
			})
		case "Commit":
			tags = append(tags, Tag{
				Name: node.Name,
				Sha:  node.Target.Commit.Oid,
				Date: node.Target.Commit.Committer.Date,
			})
		}
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, nil
}
