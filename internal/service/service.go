// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

// Package service provides interfaces for the application services.
package service

import (
	"context"
	"p86l/internal/types"
)

// InstallProgress tracks download/extract progress for a version install.
type InstallProgress struct {
	DownloadDone  int64
	DownloadTotal int64
	ExtractDone   int
	ExtractTotal  int
	Phase         types.Phase
}

// Version represents an installed game version.
type Version struct {
	Tag        string
	Executable string // absolute path
	Runnable   bool
	OS         string // "windows", "darwin", "linux"
}

// DownloadService manages downloading and extracting game versions.
type DownloadService interface {
	InstalledVersions() ([]Version, error)
	InstallVersion(ctx context.Context, url string, onProgress func(InstallProgress)) error
}

// LaunchService manages launching and killing game processes.
type LaunchService interface {
	Launch(ver Version, rebuildFn func()) error
	Kill() error
	IsRunning() bool
	IsRunningVersion(tag, os string) bool
	CurrentVersion() Version
	SetCurrentVersion(ver Version)
}

// FileOpenerService opens files, folders, and URLs in the OS.
type FileOpenerService interface {
	OpenPath(relPath string) error
	OpenURL(url string) error
	OpenFolder(folder types.Folder) error
}

// LogService provides access to application logs.
type LogService interface {
	ReadLatestLog() (string, error)
}

// ShortcutService manages desktop shortcut creation.
type ShortcutService interface {
	CreateShortcut(ver Version) error
}

// VersionService manages version deletion.
type VersionService interface {
	DeleteVersion(tag string) error
}
