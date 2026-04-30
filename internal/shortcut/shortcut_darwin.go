//go:build darwin

// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package shortcut

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func desktopDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Desktop"), nil
}

// Info.plist template for the wrapper .app bundle.
var plistTmpl = template.Must(template.New("plist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>CFBundleExecutable</key>
<string>launcher</string>
<key>CFBundleIdentifier</key>
<string>com.shortcut.{{.BundleID}}</string>
<key>CFBundleName</key>
<string>{{.Name}}</string>
<key>CFBundlePackageType</key>
<string>APPL</string>
<key>CFBundleVersion</key>
<string>1.0</string>
<key>CFBundleShortVersionString</key>
<string>1.0</string>
{{- if .Icon}}
<key>CFBundleIconFile</key>
<string>AppIcon</string>
{{- end}}
<key>LSMinimumSystemVersion</key>
<string>10.13</string>
<key>LSUIElement</key>
<false/>
<key>NSHighResolutionCapable</key>
<true/>
</dict>
</plist>
`))

// create is the macOS implementation of shortcut.Create.
func create(opts Options) (string, error) {
	if err := mustExist(opts.Target); err != nil {
		return "", err
	}

	desktop, err := desktopDir()
	if err != nil {
		return "", fmt.Errorf("shortcut: cannot find desktop: %w", err)
	}

	appPath := filepath.Join(desktop, opts.Name+".app")
	macOSDir := filepath.Join(appPath, "Contents", "MacOS")
	resDir := filepath.Join(appPath, "Contents", "Resources")

	for _, d := range []string{macOSDir, resDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return "", err
		}
	}

	// Info.plist
	bundleID := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, opts.Name)

	plistData := struct {
		Name     string
		BundleID string
		Icon     string
	}{Name: opts.Name, BundleID: bundleID, Icon: opts.Icon}

	plistFile, err := os.Create(filepath.Join(appPath, "Contents", "Info.plist"))
	if err != nil {
		return "", err
	}
	if err := plistTmpl.Execute(plistFile, plistData); err != nil {
		_ = plistFile.Close()
		return "", err
	}
	_ = plistFile.Close()

	// Launcher shell script
	workDir := filepath.Dir(opts.Target)
	target := shellQuote(opts.Target)

	launcherScript := fmt.Sprintf("#!/bin/sh\ncd %s\nexec %s\n",
		shellQuote(workDir),
		target,
	)

	launcherPath := filepath.Join(macOSDir, "launcher")
	if err := os.WriteFile(launcherPath, []byte(launcherScript), 0o755); err != nil {
		return "", fmt.Errorf("shortcut: cannot write launcher: %w", err)
	}

	// Icon (optional)
	if opts.Icon != "" {
		if err := copyFile(opts.Icon, filepath.Join(resDir, "AppIcon.icns")); err != nil {
			// Non-fatal: shortcut still works without the icon.
			fmt.Fprintf(os.Stderr, "shortcut: warning: could not copy icon: %v\n", err)
		}
	}

	return appPath, nil
}

// shellQuote wraps s in single quotes, safe for POSIX sh.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// copyFile does a simple byte-for-byte copy.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
