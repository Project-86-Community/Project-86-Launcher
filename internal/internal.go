/*
 * SPDX-License-Identifier: GPL-3.0-only
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

package internal

import (
	"eightysix/content"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
)

func IsInternetReachable() bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 204
}

func CompanyPathFolder() string {
	folderPath := content.Mgdata.ObjectPropPath("cache", "darkmode.data")
	if runtime.GOOS == "windows" {
		folderPath = strings.TrimSuffix(folderPath, "Project-86-Launcher\\cache\\darkmode.data")
	} else {
		folderPath = strings.TrimSuffix(folderPath, "Project-86-Launcher/cache/darkmode.data")
	}
	return folderPath
}

func LauncherPathFolder() string {
	folderPath := content.Mgdata.ObjectPropPath("cache", "darkmode.data")
	if runtime.GOOS == "windows" {
		folderPath = strings.TrimSuffix(folderPath, "cache\\darkmode.data")
	} else {
		folderPath = strings.TrimSuffix(folderPath, "cache/darkmode.data")
	}
	return folderPath
}

func OpenFileManager(path string) error {
	return open.Run(path)
}
