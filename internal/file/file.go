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
	"os"
	"p86l/configs"
	"p86l/internal/debug"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
)

type AppFS struct {
	GdataM *gdata.Manager
}

func (afs *AppFS) clean() string {
	colorModeFile := afs.GdataM.ObjectPropPath(configs.Data, configs.ColorModeFile)
	if runtime.GOOS == "windows" {
		return strings.TrimSuffix(colorModeFile, fmt.Sprintf("%s\\%s\\%s", configs.AppName, configs.Data, configs.ColorModeFile))
	}
	return strings.TrimSuffix(colorModeFile, fmt.Sprintf("%s/%s/%s", configs.AppName, configs.Data, configs.ColorModeFile))
}

func (afs *AppFS) OpenFileManager(appDebug *debug.Debug, path string) *debug.Error {
	log.Info().Str("Open File Manager", path).Send()
	if err := open.Run(path); err != nil {
		return appDebug.New(err, debug.FSError, debug.ErrOpenFolderFailed)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (afs *AppFS) IsDir() bool {
	if afs.GdataM.ObjectPropExists(configs.Data, configs.ColorModeFile) || afs.GdataM.ObjectPropExists(configs.Data, configs.AppScaleFile) {
		return true
	}
	return false
}

func (afs *AppFS) CompanyDir(appDebug *debug.Debug) (string, *debug.Error) {
	if afs.IsDir() {
		return afs.clean(), appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
	}

	return "", appDebug.New(errors.New("CompanyDir not found"), debug.FSError, debug.ErrDirNotFound)
}

func (afs *AppFS) LauncherDir(appDebug *debug.Debug) (string, *debug.Error) {
	if afs.IsDir() {
		return afs.clean() + configs.AppName, appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
	}

	return "", appDebug.New(errors.New("LauncherDir not found"), debug.FSError, debug.ErrDirNotFound)
}

func (afs *AppFS) LogDir(appDebug *debug.Debug) (string, *debug.Error) {
	if afs.IsDir() {
		_, err := afs.LauncherDir(appDebug)
		if err.Err != nil {
			return "", err
		}

		if runtime.GOOS == "windows" {
			return afs.clean() + configs.AppName + "\\logs", appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
		}
		return afs.clean() + configs.AppName + "/logs", appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
	}

	return "", appDebug.New(errors.New("LogDir not found"), debug.FSError, debug.ErrDirNotFound)
}

func (afs *AppFS) ClearFolder(folderPath string, appDebug *debug.Debug) *debug.Error {
	// Read all items in the directory
	items, err := os.ReadDir(folderPath)
	if err != nil {
		return appDebug.New(fmt.Errorf("failed to read directory: %w", err), debug.FSError, debug.ErrFolderClear)
	}

	// Iterate through each item and remove it
	for _, item := range items {
		itemPath := filepath.Join(folderPath, item.Name())

		// If it's a directory, remove all contents recursively
		if item.IsDir() {
			if err := os.RemoveAll(itemPath); err != nil {
				return appDebug.New(fmt.Errorf("failed to remove directory %s: %w", itemPath, err), debug.FSError, debug.ErrFolderClear)
			}
		} else {
			// If it's a file, remove it directly
			if err := os.Remove(itemPath); err != nil {
				return appDebug.New(fmt.Errorf("failed to remove file %s: %w", itemPath, err), debug.FSError, debug.ErrFolderClear)
			}
		}
	}

	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}
