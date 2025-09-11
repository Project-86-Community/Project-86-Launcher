/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Project 86 Community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"p86l/configs"
	"p86l/internal/log"
	"path/filepath"
)

// Used to make folders.
func mkdirAll(path string) error {
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("%w: %w", log.ErrMkdirAllInvalid, err)
	}
	return nil
}

type Filesystem struct {
	root *os.Root
	path string
}

func NewFilesystem(extra ...string) (*Filesystem, error) {
	var companyPath string
	if len(extra) == 1 && extra[0] != "" {
		extraPath, err := GetCompanyPath(extra[0])
		if err != nil {
			return nil, err
		}
		companyPath = extraPath
	} else {
		defaultPath, err := GetCompanyPath()
		if err != nil {
			return nil, err
		}
		companyPath = defaultPath
	}

	root, err := os.OpenRoot(companyPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrRootInvalid, err)
	}

	if err := mkdirAll(filepath.Join(companyPath, configs.AppName)); err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrMkdirAllInvalid, err)
	}

	return &Filesystem{root: root, path: companyPath}, nil
}

func (f *Filesystem) Root() *os.Root {
	return f.root
}

func (f *Filesystem) Path() string {
	return f.path
}

func (f *Filesystem) Remove(filePath string) error {
	err := f.root.Remove(filePath)
	if err != nil {
		return fmt.Errorf("%w: %w", log.ErrFileRemove, err)
	}
	return nil
}

func (f *Filesystem) Exist(filePath string) bool {
	_, err := f.root.Stat(filePath)
	return err == nil
}

func (f *Filesystem) Load(filePath string) ([]byte, error) {
	loadFile, err := f.root.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrFileLoad, err)
	}
	fileBytes, err := io.ReadAll(loadFile)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrFileLoad, err)
	}
	err = loadFile.Close()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrFileLoad, err)
	}
	return fileBytes, nil
}

func (f *Filesystem) Save(filePath string, fileBytes []byte) error {
	saveFile, err := f.root.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w: %w", log.ErrFileSave, err)
	}
	_, err = saveFile.Write(fileBytes)
	if err != nil {
		return fmt.Errorf("%w: %w", log.ErrFileSave, err)
	}
	err = saveFile.Close()
	if err != nil {
		return fmt.Errorf("%w: %w", log.ErrFileSave, err)
	}
	return nil
}

func (f *Filesystem) Close() error {
	return f.root.Close()
}
