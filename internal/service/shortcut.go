// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package service

import (
	"fmt"

	"p86l/configs"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/shortcut"

	"github.com/spf13/afero"
)

type shortcutService struct {
	afs afero.Fs
}

func NewShortcutService(afs afero.Fs) ShortcutService {
	return &shortcutService{afs: afs}
}

func (s *shortcutService) CreateShortcut(ver Version) error {
	iconRelPath := fs.IconPath()
	iconPath, err := fs.RealPath(s.afs, iconRelPath)
	if err != nil {
		logger.Warn.Printf("CreateShortcut: could not resolve icon path %s: %v", iconRelPath, err)
		iconPath = ""
	} else {
		exists, err := afero.Exists(s.afs, iconRelPath)
		if err != nil || !exists {
			logger.Warn.Printf("CreateShortcut: icon file does not exist at %s (abs: %s)", iconRelPath, iconPath)
			iconPath = ""
		}
	}

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
