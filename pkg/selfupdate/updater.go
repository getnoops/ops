package selfupdate

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Updater interface {
	GetLatest(ctx context.Context) (*ReleaseAsset, error)
	UpdateTo(ctx context.Context, asset *ReleaseAsset, exePath string) error
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

func (u *updater) UpdateTo(ctx context.Context, asset *ReleaseAsset, exePath string) error {
	reader, err := u.gh.DownloadReleaseAsset(ctx, asset.AssetId)
	if err != nil {
		return err
	}
	defer reader.Close()

	out, err := DecompressCommand(reader, asset.Filename, "ops", u.os, u.arch)
	if err != nil {
		return err
	}

	// get the directory the executable exists in
	updateDir := filepath.Dir(exePath)
	filename := filepath.Base(exePath)

	// Copy the contents of newbinary to a new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, out)
	if err != nil {
		return err
	}

	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	fp.Close()

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	if err := os.Rename(exePath, oldPath); err != nil {
		return err
	}
	// move the new executable in to become the new program
	if err := os.Rename(newPath, exePath); err != nil {
		return err
	}

	// move successful, remove the old binary if needed
	if err := os.Remove(oldPath); err != nil {
		return err
	}

	return nil
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
