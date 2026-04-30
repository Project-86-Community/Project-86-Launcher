// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package fs

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	grab "github.com/cavaliergopher/grab/v3"
	"github.com/spf13/afero"
)

var versionFromURL = regexp.MustCompile(`/download/(v[^/]+)/`)

func ParseVersion(url string) (string, error) {
	m := versionFromURL.FindStringSubmatch(url)
	if len(m) < 2 {
		return "", fmt.Errorf("could not parse version from URL: %s", url)
	}
	return m[1], nil
}

type DownloadOptions struct {
	URL      string
	Progress func(done, total int64)
}

type Downloader interface {
	Download(ctx context.Context, fs afero.Fs, opts DownloadOptions) (zipPath string, err error)
}

type GrabDownloader struct{}

func (GrabDownloader) Download(ctx context.Context, afs afero.Fs, opts DownloadOptions) (string, error) {
	version, err := ParseVersion(opts.URL)
	if err != nil {
		return "", err
	}

	destDir := filepath.Join("versions", version)
	zipName := filepath.Base(opts.URL)
	destRel := filepath.Join(destDir, zipName)

	if err := afs.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("download: mkdir: %w", err)
	}

	destAbs, err := RealPath(afs, destRel)
	if err != nil {
		return "", fmt.Errorf("download: resolve path: %w", err)
	}

	req, err := grab.NewRequest(destAbs, opts.URL)
	if err != nil {
		return "", fmt.Errorf("download: build request: %w", err)
	}
	req = req.WithContext(ctx)

	client := grab.NewClient()
	resp := client.Do(req)

	// Tick progress at 200ms intervals regardless of how fast chunks arrive.
	if opts.Progress != nil {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
	loop:
		for {
			select {
			case <-ticker.C:
				opts.Progress(resp.BytesComplete(), resp.Size())
			case <-resp.Done:
				// Final update with exact values.
				opts.Progress(resp.BytesComplete(), resp.Size())
				break loop
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	} else {
		select {
		case <-resp.Done:
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	if err := resp.Err(); err != nil {
		return "", fmt.Errorf("download: %w", err)
	}

	return destRel, nil
}

type FakeDownloader struct{}

func (FakeDownloader) Download(ctx context.Context, afs afero.Fs, opts DownloadOptions) (string, error) {
	version, err := ParseVersion(opts.URL)
	if err != nil {
		return "", err
	}
	destDir := filepath.Join("versions", version)
	zipName := filepath.Base(opts.URL)
	destRel := filepath.Join(destDir, zipName)

	if err := afs.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	// Simulate download progress over 3 seconds.
	if opts.Progress != nil {
		const total = 100 * 1024 * 1024 // fake 100MB
		const steps = 30
		for i := range steps {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(100 * time.Millisecond):
				done := int64(i+1) * (total / steps)
				opts.Progress(done, total)
			}
		}
	}

	if err := afero.WriteFile(afs, destRel, []byte("fake zip"), 0644); err != nil {
		return "", err
	}
	return destRel, nil
}
