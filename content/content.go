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

package content

import (
	"context"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/quasilyte/gdata/v2"
)

type GameFile struct {
	Tag       string
	URL       string
	Timestamp time.Time
	ExpiresIn time.Duration
}

type Changelog struct {
	Body      string
	URL       string
	Timestamp time.Time
	ExpiresIn time.Duration
}

var (
	IsInternet    = false
	Mgdata        *gdata.Manager
	GithubClient  = github.NewClient(nil)
	GithubContext = context.Background()
)

var (
	RepoOwner      = "Taliayaya"
	RepoName       = "Project-86"
	UpdateGame     = false
	DownloadStatus = -1.0
	DownloadETA    string
)
