package gitclient

import (
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Ref struct {
	Sha  string
	Name string
	Date time.Time
}

type GitClient struct {
	repo *git.Repository
}

func (gc *GitClient) GetTags() ([]*Ref, error) {
	refs, err := gc.repo.References()
	if err != nil {
		return nil, err
	}

	tags := []*Ref{}
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsTag() {
			hash, err := gc.repo.ResolveRevision(plumbing.Revision(ref.Hash().String()))
			if err != nil {
				return err
			}

			commit, err := gc.repo.CommitObject(*hash)
			if err != nil {
				return err
			}

			tags = append(tags, &Ref{
				Sha:  hash.String(),
				Name: ref.Name().Short(),
				Date: commit.Committer.When,
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})

	return tags, err
}

func (gc *GitClient) GetLogs(reverse bool) ([]*Ref, error) {
	logOpts := git.LogOptions{
		Order: git.LogOrderCommitterTime,
	}

	logsRefs, err := gc.repo.Log(&logOpts)
	if err != nil {
		return nil, err
	}

	logs := []*Ref{}
	err = logsRefs.ForEach(func(c *object.Commit) error {
		logs = append(logs, &Ref{
			Sha:  c.Hash.String(),
			Name: c.Hash.String(),
			Date: c.Committer.When,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	if reverse {
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].Date.Before(logs[j].Date)
		})
	}

	return logs, nil
}

func (gc *GitClient) GetFirstCommit() (*Ref, error) {
	logs, err := gc.GetLogs(true)
	if err != nil {
		return nil, err
	}

	return logs[0], nil
}

func NewGitClient() (*GitClient, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	// TODO: Make this configurable
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}

	return &GitClient{
		repo: repo,
	}, nil
}
