// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package fs

import (
	"fmt"
	"os"
	"p86l/assets"
	"p86l/configs"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
)

// Dir returns the OS-specific appdata directory.
func Dir() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
		if base == "" {
			return "", fmt.Errorf("%%APPDATA%% not set")
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support")
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			base = xdg
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = filepath.Join(home, ".local", "share")
		}
	}
	return filepath.Join(base, configs.Name), nil
}

// New returns a real OS filesystem jailed to the appdata dir.
// All paths are relative to the jail, escapes are prevented.
// If fake is true, the filesystem reads existing files but writes to memory,
// and never creates directories on disk.
func New(fake bool) (afero.Fs, error) {
	if fake {
		return NewFake()
	}
	root, err := Dir()
	if err != nil {
		return nil, err
	}
	base := afero.NewOsFs()
	for _, dir := range []string{"logs", "versions", "icons"} {
		if err := base.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			return nil, fmt.Errorf("fs: mkdir %s: %w", dir, err)
		}
	}
	return afero.NewBasePathFs(base, root), nil
}

// NewMem returns an in-memory filesystem with the same structure for tests.
func NewMem() afero.Fs {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("logs", 0755)
	_ = fs.MkdirAll("versions", 0755)
	_ = fs.MkdirAll("icons", 0755)
	return fs
}

// NewFake returns a filesystem that reads existing appdata files but writes to memory.
// It never creates directories or writes any file on disk.
func NewFake() (afero.Fs, error) {
	root, err := Dir()
	if err != nil {
		return nil, err
	}

	// Read‑only layer for existing files on disk
	base := afero.NewOsFs()
	jailed := afero.NewBasePathFs(base, root)
	ro := afero.NewReadOnlyFs(jailed)

	// Writable in‑memory layer for fake operations
	mem := afero.NewMemMapFs()
	_ = mem.MkdirAll("logs", 0755)
	_ = mem.MkdirAll("versions", 0755)
	_ = mem.MkdirAll("icons", 0755)

	// Copy‑on‑write: reads from ro, writes to mem
	return afero.NewCopyOnWriteFs(ro, mem), nil
}

// RealPath resolves the absolute OS path for a relative jail path.
// Only works on a real (non-memory) jailed FS.
func RealPath(fs afero.Fs, rel string) (string, error) {
	bpfs, ok := fs.(*afero.BasePathFs)
	if !ok {
		return "", fmt.Errorf("fs: RealPath requires a BasePathFs")
	}
	return bpfs.RealPath(rel)
}

// IconPath returns the relative path to the icon file for the current platform.
// The icon is stored in the "icons" directory with platform-specific extension.
func IconPath() string {
	switch runtime.GOOS {
	case "windows":
		return "icons/icon.ico"
	case "darwin":
		return "icons/icon.icns"
	default: // linux and others
		return "icons/icon.png"
	}
}

// EnsureIcons copies embedded icon files to the icons directory.
// It should be called during application initialisation.
func EnsureIcons(fs afero.Fs) error {
	// Map of icon filenames to copy
	iconFiles := []string{"icon.ico", "icon.icns", "icon.png"}

	for _, filename := range iconFiles {
		data, ok := assets.IconData[filename]
		if !ok {
			// Icon not embedded; skip
			continue
		}

		destPath := filepath.Join("icons", filename)

		// Check if file already exists
		exists, err := afero.Exists(fs, destPath)
		if err != nil {
			return fmt.Errorf("check icon %s: %w", filename, err)
		}
		if exists {
			continue
		}

		// Write icon data
		if err := afero.WriteFile(fs, destPath, data, 0644); err != nil {
			return fmt.Errorf("write icon %s: %w", filename, err)
		}
	}

	return nil
}
