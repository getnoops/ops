package selfupdate

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

type Updater interface {
	GetLatest(ctx context.Context) (*ReleaseAsset, error)
}

type updater struct {
	gh             *Github
	os             string
	arch           string
	repositorySlug string
	prerelease     bool
	draft          bool
}

// getSuffixes returns all candidates to check against the assets
func (u *updater) getSuffix() string {
	ext := ".tar.gz"
	if u.os == "windows" {
		ext = ".zip"
	}

	return fmt.Sprintf("%s-%s%s", u.os, u.arch, ext)
}

func (u *updater) GetLatest(ctx context.Context) (*ReleaseAsset, error) {
	releases, err := u.gh.ListReleases(ctx)
	if err != nil {
		return nil, err
	}

	suffix := u.getSuffix()

	for _, release := range releases {
		if release.GetDraft() && !u.draft {
			continue
		}
		if release.GetPrerelease() && !u.prerelease {
			continue
		}
		if u.prerelease && !release.GetPrerelease() {
			continue
		}
		if u.draft && !release.GetDraft() {
			continue
		}

		assets := release.Assets
		if len(assets) == 0 {
			continue
		}

		// what about the suffix?
		for _, asset := range assets {
			name := strings.ToLower(asset.GetName())
			if strings.HasSuffix(name, suffix) {
				return &ReleaseAsset{
					TagName:  release.GetTagName(),
					Filename: asset.GetName(),
					AssetId:  asset.GetID(),
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no release found")
}

func NewUpdater(repositorySlug string, prerelease bool, draft bool) (Updater, error) {
	gh, err := NewGithub(repositorySlug)
	if err != nil {
		return nil, err
	}

	return &updater{
		gh:             gh,
		os:             runtime.GOOS,
		arch:           runtime.GOARCH,
		repositorySlug: repositorySlug,
		prerelease:     prerelease,
		draft:          draft,
	}, nil
}
