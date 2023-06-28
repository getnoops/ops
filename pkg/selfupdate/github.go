package selfupdate

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v53/github"
)

type Github struct {
	api   *github.Client
	owner string
	repo  string
}

func (g *Github) ListReleases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	releases, res, err := g.api.Repositories.ListReleases(ctx, g.owner, g.repo, nil)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return nil, errors.New("repository not found")
		}
		return nil, err
	}
	return releases, nil
}

func (g *Github) DownloadReleaseAsset(ctx context.Context, assetId int64) (io.ReadCloser, error) {
	client := http.DefaultClient
	asset, _, err := g.api.Repositories.DownloadReleaseAsset(ctx, g.owner, g.repo, assetId, client)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func NewGithub(repositorySlug string) (*Github, error) {
	out := strings.Split(repositorySlug, "/")
	if len(out) != 2 {
		return nil, errors.New("invalid repository slug")
	}

	owner := out[0]
	repo := out[1]

	return &Github{
		api: github.NewClient(nil),

		owner: owner,
		repo:  repo,
	}, nil
}
