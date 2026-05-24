// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package service

import (
	"os"
	"path/filepath"

	"p86l/internal/fs"
	"p86l/internal/logger"

	"github.com/spf13/afero"
)

type versionService struct {
	afs afero.Fs
}

func NewVersionService(afs afero.Fs) VersionService {
	return &versionService{afs: afs}
}

func (s *versionService) DeleteVersion(tag string) error {
	path, err := fs.RealPath(s.afs, filepath.Join("versions", tag))
	if err != nil {
		return err
	}
	logger.Info.Printf("deleting version: %s", tag)
	return os.RemoveAll(path)
}
