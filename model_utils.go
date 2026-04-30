// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package p86l

import (
	"fmt"
	"os"
	"p86l/configs"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/shortcut"
	"p86l/internal/types"
	"path/filepath"
	"sort"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
)

func (m *Model) ReadLatestLog() (string, error) {
	if m.fake {
		return "[fake mode - no real log file]", nil
	}

	logsDir, err := fs.RealPath(m.fs, "logs")
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return "", err
	}

	var logs []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, configs.LogPrefix) && strings.HasSuffix(name, configs.LogExt) {
			logs = append(logs, name)
		}
	}

	if len(logs) == 0 {
		return "[no previous session log found]", nil
	}

	sort.Strings(logs)
	prev := filepath.Join(logsDir, logs[len(logs)-1])

	data, err := os.ReadFile(prev)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Model) Open(input string, isUrl bool) {
	if m.fake {
		return
	}

	logger.Info.Printf("opening: %s", input)

	if isUrl {
		_ = open.Start(input)
		return
	}

	path, err := fs.RealPath(m.fs, filepath.Join("versions", input))
	if err != nil {
		return
	}
	_ = open.Start(path)
}

func (m *Model) OpenFolder(folder types.Folder) {
	if m.fake {
		return
	}

	var rel string
	switch folder {
	case types.FolderRoot:
		rel = "."
	case types.FolderVersions:
		rel = "versions"
	case types.FolderLogs:
		rel = "logs"
	default:
		return
	}

	path, err := fs.RealPath(m.fs, rel)
	if err != nil {
		logger.Warn.Printf("OpenFolder: could not resolve %q: %v", rel, err)
		return
	}

	logger.Info.Printf("opening folder in file manager: %s", rel)
	_ = open.Start(path)
}

func (m *Model) DeleteVersion(tag string) error {
	if m.fake {
		logger.Info.Printf("DeleteVersion: fake mode, skipping delete for %s", tag)
		return nil
	}

	path, err := fs.RealPath(m.fs, filepath.Join("versions", tag))
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func (m *Model) CreateShortcut(ver Version) error {
	if m.fake {
		logger.Info.Printf("CreateShortcut: fake mode, skipping shortcut creation for %s", ver.Tag)
		return nil
	}

	// Get the absolute path to the icon file
	iconRelPath := fs.IconPath()
	iconPath, err := fs.RealPath(m.fs, iconRelPath)
	if err != nil {
		logger.Warn.Printf("CreateShortcut: could not resolve icon path %s: %v", iconRelPath, err)
		// Continue without icon
		iconPath = ""
	} else {
		// Check if icon file exists using the filesystem
		exists, err := afero.Exists(m.fs, iconRelPath)
		if err != nil || !exists {
			logger.Warn.Printf("CreateShortcut: icon file does not exist at %s (abs: %s)", iconRelPath, iconPath)
			iconPath = ""
		}
	}

	// Create shortcut name in format "Project 86 v1.0.0-alpha"
	shortcutName := fmt.Sprintf("%s %s", configs.Game, ver.Tag)

	opts := shortcut.Options{
		Name:   shortcutName,
		Target: ver.Executable,
		Icon:   iconPath,
	}

	logger.Info.Printf("Creating desktop shortcut: %s -> %s (icon: %s)",
		shortcutName, ver.Executable, iconPath)

	createdPath, err := shortcut.Create(opts)
	if err != nil {
		return fmt.Errorf("CreateShortcut: %w", err)
	}

	logger.Info.Printf("Shortcut created successfully at: %s", createdPath)
	return nil
}
