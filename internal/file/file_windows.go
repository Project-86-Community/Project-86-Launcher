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
	"p86l/internal/log"
	"path/filepath"
)

func GetCompanyPath(extra ...string) (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", log.ErrCompanyPathAppData
	}
	companyPath := filepath.Join(appData, configs.CompanyName)
	// Used for testing only!
	if len(extra) == 1 && extra[0] != "" {
		companyPath = fmt.Sprintf("%s_%s", companyPath, extra[0])
	}
	if err := mkdirAll(companyPath); err != nil {
		return "", err
	}
	return companyPath, nil
}
