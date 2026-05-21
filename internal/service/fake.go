// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package service

import (
	"context"
	"fmt"

	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"

	"github.com/spf13/afero"
)

// FakeDownloadService returns canned versions and simulates installs.
type FakeDownloadService struct {
	afs       afero.Fs
	dl        fs.Downloader
	ex        fs.Extractor
	FakeError bool
}

func NewFakeDownloadService(afs afero.Fs) *FakeDownloadService {
	return &FakeDownloadService{
		afs: afs,
		dl:  fs.FakeDownloader{},
		ex:  fs.FakeExtractor{},
	}
}

func (s *FakeDownloadService) InstalledVersions() ([]Version, error) {
	return []Version{
		{Tag: "v0.0.0-alpha", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
		{Tag: "v1.0.0-beta", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
	}, nil
}

func (s *FakeDownloadService) InstallVersion(ctx context.Context, url string, onProgress func(InstallProgress)) error {
	if s.FakeError {
		return fmt.Errorf("install failed: Loren ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
	}

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
		return err
	}

	p.Phase = types.PhaseExtract
	if onProgress != nil {
		onProgress(p)
	}

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
		return err
	}

	p.Phase = types.PhaseDone
	if onProgress != nil {
		onProgress(p)
	}
	return nil
}

// FakeLaunchService simulates process launch/kill.
type FakeLaunchService struct {
	version    Version
	isRunning  bool
	runningTag string
	runningOS  string
}

func NewFakeLaunchService() *FakeLaunchService {
	return &FakeLaunchService{}
}

func (s *FakeLaunchService) CurrentVersion() Version       { return s.version }
func (s *FakeLaunchService) SetCurrentVersion(ver Version) { s.version = ver }
func (s *FakeLaunchService) IsRunningVersion(tag, os string) bool {
	return s.isRunning && s.runningTag == tag && s.runningOS == os
}
func (s *FakeLaunchService) IsRunning() bool { return s.isRunning }

func (s *FakeLaunchService) Launch(ver Version, rebuildFn func()) error {
	logger.Info.Printf("Launch: fake mode, skipping launch for %s", ver.Tag)
	s.isRunning = true
	s.runningTag = ver.Tag
	s.runningOS = ver.OS
	return nil
}

func (s *FakeLaunchService) Kill() error {
	logger.Info.Printf("Kill: fake mode, skipping kill for %s", s.version.Tag)
	s.isRunning = false
	s.runningTag = ""
	s.runningOS = ""
	return nil
}

// FakeFileOpenerService no-ops all open operations.
type FakeFileOpenerService struct{}

func NewFakeFileOpenerService() *FakeFileOpenerService                { return &FakeFileOpenerService{} }
func (s *FakeFileOpenerService) OpenPath(relPath string) error        { return nil }
func (s *FakeFileOpenerService) OpenURL(url string) error             { return nil }
func (s *FakeFileOpenerService) OpenFolder(folder types.Folder) error { return nil }

// FakeLogService returns a canned log message.
type FakeLogService struct{}

func NewFakeLogService() *FakeLogService { return &FakeLogService{} }
func (s *FakeLogService) ReadLatestLog() (string, error) {
	return "[fake mode - no real log file]", nil
}

// FakeShortcutService no-ops shortcut creation.
type FakeShortcutService struct{}

func NewFakeShortcutService() *FakeShortcutService { return &FakeShortcutService{} }
func (s *FakeShortcutService) CreateShortcut(ver Version) error {
	logger.Info.Printf("CreateShortcut: fake mode, skipping shortcut creation for %s", ver.Tag)
	return nil
}

// FakeVersionService no-ops version deletion.
type FakeVersionService struct{}

func NewFakeVersionService() *FakeVersionService { return &FakeVersionService{} }
func (s *FakeVersionService) DeleteVersion(tag string) error {
	logger.Info.Printf("DeleteVersion: fake mode, skipping delete for %s", tag)
	return nil
}
