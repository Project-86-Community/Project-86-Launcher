//go:build windows

// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package main

import (
	"fmt"
	"os"
	"p86l/configs"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func tryLock() (*os.File, error) {
	lockPath := filepath.Join(os.TempDir(), configs.LockFile)
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	ol := new(windows.Overlapped)
	err = windows.LockFileEx(
		windows.Handle(f.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0, 1, 0, ol,
	)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("another instance is already running")
	}

	return f, nil
}
