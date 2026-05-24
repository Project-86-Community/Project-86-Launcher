//go:build linux

// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package shortcut

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func desktopDir() (string, error) {
	// Respect XDG_DESKTOP_DIR if set, fall back to ~/Desktop.
	if xdg := os.Getenv("XDG_DESKTOP_DIR"); xdg != "" {
		return xdg, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Desktop"), nil
}

var desktopTmpl = template.Must(template.New("desktop").Parse(
	`[Desktop Entry]
	Version=1.0
	Type=Application
	Name={{.Name}}
	Exec={{.Exec}}
	{{- if .Icon}}
	Icon={{.Icon}}
	{{- end}}
	Path={{.WorkingDir}}
	Terminal=false
	StartupNotify=true
	`))

// create is the Linux implementation of shortcut.Create.
func create(opts Options) (string, error) {
	if err := mustExist(opts.Target); err != nil {
		return "", err
	}

	desktop, err := desktopDir()
	if err != nil {
		return "", fmt.Errorf("shortcut: cannot find desktop: %w", err)
	}
	if err := os.MkdirAll(desktop, 0o755); err != nil {
		return "", err
	}

	destPath := filepath.Join(desktop, sanitizeFilename(opts.Name)+".desktop")

	workDir := filepath.Dir(opts.Target)

	// Build the Exec value: just the target, shell-quoted.
	exec := quoteArg(opts.Target)

	data := struct {
		Name       string
		Exec       string
		Icon       string
		WorkingDir string
	}{
		Name:       opts.Name,
		Exec:       exec,
		Icon:       opts.Icon,
		WorkingDir: workDir,
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return "", fmt.Errorf("shortcut: cannot create file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := desktopTmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("shortcut: template error: %w", err)
	}

	return destPath, nil
}

// quoteArg wraps a string in single quotes, escaping any single quotes inside.
func quoteArg(s string) string {
	if !strings.ContainsAny(s, " \t\n'\"\\") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// sanitizeFilename replaces characters that are awkward in file names.
func sanitizeFilename(name string) string {
	r := strings.NewReplacer("/", "-", "\\", "-", ":", "-")
	return r.Replace(name)
}
