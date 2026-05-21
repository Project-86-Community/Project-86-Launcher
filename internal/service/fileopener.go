// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package service

import (
	"path/filepath"

	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
)

type fileOpenerService struct {
	afs afero.Fs
}

func NewFileOpenerService(afs afero.Fs) FileOpenerService {
	return &fileOpenerService{afs: afs}
}

func (s *fileOpenerService) OpenPath(relPath string) error {
	path, err := fs.RealPath(s.afs, filepath.Join("versions", relPath))
	if err != nil {
		return err
	}
	logger.Info.Printf("opening: %s", path)
	return open.Start(path)
}

func (s *fileOpenerService) OpenURL(url string) error {
	logger.Info.Printf("opening URL: %s", url)
	return open.Start(url)
}

func (s *fileOpenerService) OpenFolder(folder types.Folder) error {
	var rel string
	switch folder {
	case types.FolderRoot:
		rel = "."
	case types.FolderVersions:
		rel = "versions"
	case types.FolderLogs:
		rel = "logs"
	default:
		return nil
	}

	path, err := fs.RealPath(s.afs, rel)
	if err != nil {
		logger.Warn.Printf("OpenFolder: could not resolve %q: %v", rel, err)
		return err
	}

	logger.Info.Printf("opening folder in file manager: %s", rel)
	return open.Start(path)
}
