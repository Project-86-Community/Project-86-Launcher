//go:build darwin || linux

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
	"fmt"
	"os"
	"p86l/configs"
	"p86l/internal/debug"
	"path/filepath"
)

func GetCompanyPath(appDebug *debug.Debug, extra ...string) (string, *debug.Error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", appDebug.New(err, debug.FSError, debug.ErrFSDirInvalid)
	}
	dataPath := filepath.Join(home, ".local", "share", configs.CompanyName)
	// Used for testing only!
	if len(extra) == 1 && extra[0] != "" {
		dataPath = fmt.Sprintf("%s_%s", dataPath, extra[0])
	}
	if dErr := mkdirAll(appDebug, dataPath); dErr != nil {
		return "", dErr
	}
	return dataPath, nil
}
