//go:build windows

// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package shortcut

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func desktopDir() (string, error) {
	// USERPROFILE is set on all modern Windows versions.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Desktop"), nil
}

// create is the Windows implementation of shortcut.Create.
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

	lnkPath := filepath.Join(desktop, opts.Name+".lnk")

	workDir := filepath.Dir(opts.Target)

	// Build a PowerShell script that creates the shortcut via COM.
	ps := fmt.Sprintf(`
	$ws  = New-Object -ComObject WScript.Shell
	$sc  = $ws.CreateShortcut('%s')
	$sc.TargetPath       = '%s'
	$sc.WorkingDirectory = '%s'
	%s
	$sc.Save()
	`,
		escapePSString(lnkPath),
		escapePSString(opts.Target),
		escapePSString(workDir),
		iconLine(opts.Icon),
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("shortcut: PowerShell failed: %w", err)
	}

	return lnkPath, nil
}

// iconLine returns the PowerShell line that sets the icon, or empty string.
func iconLine(icon string) string {
	if icon == "" {
		return ""
	}
	return fmt.Sprintf("$sc.IconLocation = '%s'", escapePSString(icon))
}

// escapePSString escapes single quotes for embedding in a PS single-quoted string.
func escapePSString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
