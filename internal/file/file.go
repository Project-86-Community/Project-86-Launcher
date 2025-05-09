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
	"p86l/internal/debug"
	"path/filepath"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/skratchdot/open-golang/open"
)

func mkdirAll(appDebug *debug.Debug, path string) *debug.Error {
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrFSDirNew)
	}
	return nil
}

type Data struct {
	Locale    string
	AppScale  int
	ColorMode guigui.ColorMode
}

type Cache struct {
	Repo      *github.RepositoryRelease
	Timestamp time.Time
	ExpiresIn time.Duration
}

func (c *Cache) Validate(appDebug *debug.Debug) *debug.Error {
	if c.Repo == nil {
		if c.Repo.GetBody() == "" {
			if c.Repo.GetHTMLURL() == "" {
				if c.Repo.Assets == nil {
					return appDebug.New(errors.New("Assets is empty"), debug.CacheError, debug.ErrCacheAssetsInvalid)
				}
				return appDebug.New(errors.New("URL is empty"), debug.CacheError, debug.ErrCacheURLInvalid)
			}
			return appDebug.New(errors.New("Body is empty"), debug.CacheError, debug.ErrCacheBodyInvalid)
		}
		return appDebug.New(errors.New("Repo is empty"), debug.CacheError, debug.ErrCacheRepoInvalid)
	}

	return nil
}

type AppFS struct {
	Root           *os.Root
	CompanyDirPath string
}

func NewFS(appDebug *debug.Debug, extra ...string) (*AppFS, *debug.Error) {
	var companyPath string
	if len(extra) == 1 && extra[0] != "" {
		cPath, dErr := GetCompanyPath(appDebug, extra[0])
		if dErr != nil {
			return nil, dErr
		}
		companyPath = cPath
	} else {
		cPath, dErr := GetCompanyPath(appDebug)
		if dErr != nil {
			return nil, dErr
		}
		companyPath = cPath
	}

	dErr := mkdirAll(appDebug, filepath.Join(companyPath, configs.AppName))
	if dErr != nil {
		return nil, dErr
	}

	root, err := os.OpenRoot(filepath.Join(companyPath, configs.AppName))
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFSRootInvalid)
	}

	return &AppFS{
		Root:           root,
		CompanyDirPath: companyPath,
	}, nil
}

func (a *AppFS) OpenFileManager(appDebug *debug.Debug, path string) *debug.Error {
	if err := open.Run(path); err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrFSOpenFileManagerInvalid)
	}
	return nil
}

func (a *AppFS) IsDir(appDebug *debug.Debug, filePath string) *debug.Error {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return appDebug.New(err, debug.FSError, debug.ErrFSRootFileNotExist)
		}
		return appDebug.New(err, debug.FSError, debug.ErrFSRootFileInvalid)
	}
	return nil
}

func (a *AppFS) Stat(appDebug *debug.Debug, statFile string) *debug.Error {
	_, err := a.Root.Stat(statFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return appDebug.New(err, debug.FSError, debug.ErrFSRootFileNotExist)
		}
		return appDebug.New(err, debug.FSError, debug.ErrFSRootFileInvalid)
	}
	return nil
}

func (a *AppFS) Save(appDebug *debug.Debug, saveFile string, bytes []byte) *debug.Error {
	file, err := a.Root.Create(saveFile)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrFSRootFileNew)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrFSRootFileWrite)
	}

	return nil
}

func (a *AppFS) Load(appDebug *debug.Debug, loadFile string) ([]byte, *debug.Error) {
	dErr := a.Stat(appDebug, loadFile)
	if dErr != nil {
		return nil, dErr
	}

	file, err := a.Root.Open(loadFile)
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFSRootFileInvalid)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFSRootFileRead)
	}

	return b, nil
}

// Project-86-Community/Project-86-Launcher
func (a *AppFS) DirAppPath() string {
	return filepath.Join(a.CompanyDirPath, configs.AppName)
}

// Project-86-Community/Project-86-Launcher/data.gob
func (a *AppFS) FileDataPath() string {
	return filepath.Join(a.DirAppPath(), configs.DataFile)
}

// Project-86-Community/Project-86-Launcher/cache.gob
func (a *AppFS) FileCachePath() string {
	return filepath.Join(a.DirAppPath(), configs.CacheFile)
}
