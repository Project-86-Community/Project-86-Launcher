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
	"io/fs"
	"os"
	"p86l/configs"
	"p86l/internal/debug"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
)

func mkdirAll(appDebug *debug.Debug, path string) *debug.Error {
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrNewDirFailed)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

type AppFS struct {
	Root           *os.Root
	CompanyDirPath string
}

func NewFS(appDebug *debug.Debug, root *os.Root, companyPath, appDirPath string) (*AppFS, *debug.Error) {
	if _, err := root.Stat(appDirPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := root.Mkdir(appDirPath, 0755)
			if err != nil {
				return nil, appDebug.New(err, debug.FSError, debug.ErrFSOpenFailed)
			}
		}
	}

	return &AppFS{
		Root:           root,
		CompanyDirPath: companyPath,
	}, appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (a *AppFS) OpenFileManager(appDebug *debug.Debug, path string) *debug.Error {
	if err := open.Run(path); err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrOpenFolderFailed)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (a *AppFS) IsDir(appDebug *debug.Debug, path string) *debug.Error {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) && err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrDirNotFound)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (a *AppFS) Save(appDebug *debug.Debug, key, saveFile string, bytes []byte) *debug.Error {
	savePath := filepath.Join(a.CompanyDirPath, configs.AppName, key)
	if _, err := os.Stat(savePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.Mkdir(savePath, 0755)
			if err != nil {
				return appDebug.New(err, debug.FSError, debug.ErrNewFileFailed)
			}
		}
	}
	root, err := os.OpenRoot(savePath)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrFileNotFound)
	}
	file, err := root.Create(saveFile)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrOpenFolderFailed)
	}
	_, err = file.Write(bytes)
	if err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrNewFileFailed)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (a *AppFS) Load(appDebug *debug.Debug, key, loadFile string) ([]byte, *debug.Error) {
	loadPath := filepath.Join(a.CompanyDirPath, configs.AppName, key)
	root, err := os.OpenRoot(loadPath)
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrOpenFolderFailed)
	}
	if _, err := root.Stat(loadFile); err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFileNotFound)
	}
	bytes, err := os.ReadFile(filepath.Join(loadPath, loadFile))
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFileNotFound)
	}
	return bytes, appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (a *AppFS) LogDir(appDebug *debug.Debug) (string, *debug.Error) {
	logsPath := filepath.Join(a.CompanyDirPath, configs.AppName, "logs")
	if dErr := a.IsDir(appDebug, logsPath); dErr.Err != nil {
		return "", dErr
	}
	return logsPath, appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

//
// func (a *AppFS) ResetData(appDebug *debug.Debug) *debug.Error {
// 	if dErr := a.IsDir(appDebug, filepath.Join(a.AppDirPath, configs.Data)); dErr.Err != nil {
// 		return dErr
// 	}
// 	if err := os.RemoveAll(filepath.Join(a.AppDirPath, configs.Data)); err != nil {
// 		return appDebug.New(err, debug.FSError, debug.ErrFolderClear)
// 	}
// 	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
// }
//
// func (a *AppFS) ResetCache(appDebug *debug.Debug) *debug.Error {
// 	if dErr := a.IsDir(appDebug, filepath.Join(a.AppDirPath, configs.Cache)); dErr.Err != nil {
// 		return dErr
// 	}
// 	if err := os.RemoveAll(filepath.Join(a.AppDirPath, configs.Cache)); err != nil {
// 		return appDebug.New(err, debug.FSError, debug.ErrFolderClear)
// 	}
// 	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
// }
//
// func (a *AppFS) ResetAll(appDebug *debug.Debug) *debug.Error {
// 	if err := a.ResetData(appDebug); err.Err != nil {
// 		return err
// 	}
// 	if err := a.ResetCache(appDebug); err.Err != nil {
// 		return err
// 	}
// 	if dErr := a.IsDir(appDebug, filepath.Join(a.AppDirPath, "logs")); dErr.Err != nil {
// 		return dErr
// 	}
// 	if err := os.RemoveAll(filepath.Join(a.AppDirPath, "logs")); err != nil {
// 		return appDebug.New(err, debug.FSError, debug.ErrFolderClear)
// 	}
//
// 	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
// }
