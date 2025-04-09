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

package app

import (
	"net/http"
	"p86l/internal/cache"
	"p86l/internal/data"
	"p86l/internal/debug"
	"p86l/internal/file"
	"time"
)

type App struct {
	isInternet bool

	Debug *debug.Debug
	FS    *file.AppFS
	Data  *data.Data
	Cache *cache.Cache
}

func (a *App) IsInternet() bool {
	return a.isInternet
}

func (a *App) isInternetReachable() bool {
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

func (a *App) UpdateInternet() {
	if a.isInternetReachable() {
		a.isInternet = true
	} else {
		a.isInternet = false
	}
}
