//go:build !windows

package main

import (
	"fmt"
	"os"
	"p86l/configs"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func tryLock() (*os.File, error) {
	lockPath := filepath.Join(os.TempDir(), configs.LockFile)
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	err = unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("another instance is already running")
	}

	return f, nil
}
