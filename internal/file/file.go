/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Ilan Mayeux
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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
	"eightysix/configs"
	"fmt"
	"runtime"
	"strings"

	"github.com/quasilyte/gdata/v2"
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

func (afs *AppFS) OpenFileManager(path string) error {
	return open.Run(path)
}

func (afs *AppFS) IsDir() bool {
	if afs.GdataM.ObjectPropExists(configs.Data, configs.ColorModeFile) || afs.GdataM.ObjectPropExists(configs.Data, configs.AppScaleFile) {
		return true
	}
	return false
}

func (afs *AppFS) CompanyDir() (string, error) {
	if afs.IsDir() {
		return afs.clean(), nil
	}

	return "", nil
}

func (afs *AppFS) LauncherDir() (string, error) {
	if afs.IsDir() {
		return afs.clean() + configs.AppName, nil
	}

	return "", nil
}
