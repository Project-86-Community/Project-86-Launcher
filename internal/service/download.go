// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"

	"github.com/spf13/afero"
)

type downloadService struct {
	afs afero.Fs
	dl  fs.Downloader
	ex  fs.Extractor
}

func NewDownloadService(afs afero.Fs, dl fs.Downloader, ex fs.Extractor) DownloadService {
	return &downloadService{afs: afs, dl: dl, ex: ex}
}

// executable candidates across all platforms.
var exeCandidates = []struct {
	path string
	os   string
}{
	{"Project-86.exe", "windows"},
	{"Project86.exe", "windows"},
	{"Project-86.x86_64", "linux"},
	{"Project86.x86_64", "linux"},
	{"Project-86.x64", "linux"},
	{"Project86.x64", "linux"},
	{"Project-86", "linux"},
	{"Project86", "linux"},
}

func (s *downloadService) InstalledVersions() ([]Version, error) {
	versionsDir, err := fs.RealPath(s.afs, "versions")
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	var versions []Version
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		tag := entry.Name()
		tagDir := filepath.Join(versionsDir, tag)

		subEntries, err := os.ReadDir(tagDir)
		if err != nil {
			continue
		}

		for _, sub := range subEntries {
			if !sub.IsDir() {
				continue
			}
			gameDir := filepath.Join(tagDir, sub.Name())

			// Check for macOS .app bundles first.
			if vers := findMacOSExecutables(gameDir, tag); len(vers) > 0 {
				versions = append(versions, vers...)
				break
			}

			// Check for regular executables (Windows/Linux).
			if ver, ok := findExe(gameDir, tag); ok {
				versions = append(versions, ver)
				break
			}
		}
	}

	return versions, nil
}

func findMacOSExecutables(gameDir, tag string) []Version {
	subEntries, _ := os.ReadDir(gameDir)
	var versions []Version
	for _, sse := range subEntries {
		if !sse.IsDir() || !strings.HasSuffix(sse.Name(), ".app") {
			continue
		}
		appDir := filepath.Join(gameDir, sse.Name())
		macOSDir := filepath.Join(appDir, "Contents", "MacOS")
		macContents, _ := os.ReadDir(macOSDir)
		for _, mc := range macContents {
			if mc.IsDir() {
				continue
			}
			versions = append(versions, Version{
				Tag:        tag,
				Executable: filepath.Join(macOSDir, mc.Name()),
				Runnable:   runtime.GOOS == "darwin",
				OS:         "darwin",
			})
		}
	}
	return versions
}

func findExe(gameDir, tag string) (Version, bool) {
	for _, cand := range exeCandidates {
		candidate := filepath.Join(gameDir, cand.path)
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return Version{
				Tag:        tag,
				Executable: candidate,
				Runnable:   runtime.GOOS == cand.os,
				OS:         cand.os,
			}, true
		}
	}
	return Version{}, false
}

func (s *downloadService) InstallVersion(ctx context.Context, url string, onProgress func(InstallProgress)) error {
	version, err := fs.ParseVersion(url)
	if err != nil {
		return fmt.Errorf("InstallVersion: could not parse version from URL: %w", err)
	}
	logger.Info.Printf("install started - version: %s url: %s", version, url)

	p := InstallProgress{Phase: types.PhaseDownload}

	zipPath, err := s.dl.Download(ctx, s.afs, fs.DownloadOptions{
		URL: url,
		Progress: func(done, total int64) {
			p.DownloadDone = done
			p.DownloadTotal = total
			if onProgress != nil {
				onProgress(p)
			}
		},
	})
	if err != nil {
		logger.Error.Printf("download failed [%s]: %v", version, err)
		return fmt.Errorf("InstallVersion: download: %w", err)
	}
	logger.Info.Printf("download complete [%s]: %s", version, zipPath)

	p.Phase = types.PhaseExtract
	if onProgress != nil {
		onProgress(p)
	}
	logger.Info.Printf("extraction started [%s]", version)

	err = s.ex.Extract(ctx, s.afs, zipPath, fs.ExtractOptions{
		OnFile: func(extracted, total int) {
			p.ExtractDone = extracted
			p.ExtractTotal = total
			if onProgress != nil {
				onProgress(p)
			}
		},
	})
	if err != nil {
		logger.Error.Printf("extraction failed [%s]: %v", version, err)
		return fmt.Errorf("InstallVersion: extract: %w", err)
	}
	logger.Info.Printf("extraction complete [%s]", version)
	logger.Info.Printf("install finished [%s]", version)

	p.Phase = types.PhaseDone
	if onProgress != nil {
		onProgress(p)
	}
	return nil
}
