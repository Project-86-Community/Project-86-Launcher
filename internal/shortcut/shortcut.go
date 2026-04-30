// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only
//
// Package shortcut creates native desktop shortcuts for executables.
// Supports Windows (.lnk), Linux (.desktop), and macOS (.app alias).

package shortcut

import (
	"fmt"
	"os"
)

// Options holds the configuration for a desktop shortcut.
type Options struct {
	// Name is the display name shown on the desktop.
	Name string

	// Target is the absolute path to the executable.
	Target string

	// Icon is the path to an icon file (optional).
	// Windows: .ico file
	// Linux:   .png / .svg / icon name from theme
	// macOS:   .icns file
	Icon string
}

// Create places a desktop shortcut for the executable described by opts.
// The shortcut is written to the user's desktop directory.
// Returns the full path of the created shortcut file or an error.
func Create(opts Options) (string, error) {
	if opts.Name == "" {
		return "", fmt.Errorf("shortcut: Name is required")
	}
	if opts.Target == "" {
		return "", fmt.Errorf("shortcut: Target is required")
	}
	return create(opts) // implemented per-platform
}

// DesktopDir returns the path to the current user's desktop directory.
func DesktopDir() (string, error) {
	return desktopDir()
}

// mustExist returns an error if path does not exist.
func mustExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("shortcut: path does not exist: %s", path)
	}
	return nil
}
