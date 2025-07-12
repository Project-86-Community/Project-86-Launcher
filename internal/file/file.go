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
	"io"
	"io/fs"
	"os"
	"p86l/configs"
	pd "p86l/internal/debug"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
)

// Used to make folders.
func mkdirAll(appDebug *pd.Debug, path string) *pd.Error {
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return appDebug.New(err, pd.FSError, pd.ErrFSDirNew)
	}
	return nil
}

type AppFS struct {
	Root           *os.Root
	CompanyDirPath string
}

// Make new FS for app.
func NewFS(appDebug *pd.Debug, extra ...string) (*AppFS, *pd.Error) {
	// Handles the company path and a path for debugging.
	var companyPath string
	if len(extra) == 1 && extra[0] != "" {
		cPath, err := GetCompanyPath(appDebug, extra[0])
		if err != nil {
			return nil, err
		}
		companyPath = cPath
	} else {
		cPath, err := GetCompanyPath(appDebug)
		if err != nil {
			return nil, err
		}
		companyPath = cPath
	}

	// Makes the path for company and app.
	err := mkdirAll(appDebug, filepath.Join(companyPath, configs.AppName))
	if err != nil {
		return nil, err
	}

	// Creates a virtual filesystem thats in company path, that protects/restricts changes outside of it.
	root, rErr := os.OpenRoot(companyPath)
	if rErr != nil {
		return nil, appDebug.New(rErr, pd.FSError, pd.ErrFSRootInvalid)
	}

	return &AppFS{
		Root:           root,
		CompanyDirPath: companyPath,
	}, nil
}

// Opens the filemanager app with the given path.
func (a *AppFS) OpenFileManager(appDebug *pd.Debug, path string) *pd.Error {
	if err := open.Run(path); err != nil {
		return appDebug.New(err, pd.FSError, pd.ErrFSOpenFileManagerInvalid)
	}
	log.Info().Str("Path", path).Str("AppFS", "OpenFileManager").Msg("FileManager")
	return nil
}

// Checks if the directory exists, uses OS
func (a *AppFS) IsDir(appDebug *pd.Debug, filePath string) *pd.Error {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return appDebug.New(err, pd.FSError, pd.ErrFSRootFileNotExist)
		}
		return appDebug.New(err, pd.FSError, pd.ErrFSRootFileInvalid)
	}
	return nil
}

// same as `IsDir` but its restricted via os.Root
func (a *AppFS) IsDirR(appDebug *pd.Debug, statFile string) *pd.Error {
	_, err := a.Root.Stat(statFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return appDebug.New(err, pd.FSError, pd.ErrFSRootFileNotExist)
		}
		return appDebug.New(err, pd.FSError, pd.ErrFSRootFileInvalid)
	}
	return nil
}

// Saves a file to disk.
func (a *AppFS) Save(appDebug *pd.Debug, saveFile string, bytes []byte) *pd.Error {
	file, err := a.Root.Create(saveFile)
	if err != nil {
		return appDebug.New(err, pd.FSError, pd.ErrFSRootFileNew)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return appDebug.New(err, pd.FSError, pd.ErrFSRootFileWrite)
	}

	return nil
}

// Loads binary data from a file in disk.
func (a *AppFS) Load(appDebug *pd.Debug, loadFile string) ([]byte, *pd.Error) {
	dErr := a.IsDirR(appDebug, loadFile)
	if dErr != nil {
		return nil, dErr
	}

	file, err := a.Root.Open(loadFile)
	if err != nil {
		return nil, appDebug.New(err, pd.FSError, pd.ErrFSRootFileInvalid)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, appDebug.New(err, pd.FSError, pd.ErrFSRootFileRead)
	}

	return b, nil
}

// /Project-86-Launcher/
func (a *AppFS) DirAppPath() string {
	return configs.AppName
}

// /build/
func (a *AppFS) DirBuildPath() string {
	return "build"
}

// /build/game
func (a *AppFS) DirGamePath() string {
	return filepath.Join(a.DirBuildPath(), "game")
}

// /build/prerelease
func (a *AppFS) DirPreReleasePath() string {
	return filepath.Join(a.DirBuildPath(), "prerelease")
}

// -- files --

// /Project-86-Launcher/data.json
func (a *AppFS) FileDataPath() string {
	return filepath.Join(a.DirAppPath(), configs.DataFile)
}

// /Project-86-Launcher/cache.json
func (a *AppFS) FileCachePath() string {
	return filepath.Join(a.DirAppPath(), configs.CacheFile)
}
