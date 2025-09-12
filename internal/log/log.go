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

package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var ErrLogFileInvalid = errors.New("failed to create log file")

type Manager int

const (
	UnknownManager Manager = iota
	AppManager
	ErrorManager
	FileManager
	NetworkManager
)

func (m Manager) String() string {
	list := []string{"Unknown", "App", "Error", "File", "Network"}
	return list[m] + "Manager"
}

// -- errors --

var (
	ErrMkdirAllInvalid    = errors.New("failed to create new folder")
	ErrCompanyPathAppData = errors.New("failed to get appdata")
	ErrRootInvalid        = errors.New("failed to open root")

	ErrFileRemove = errors.New("failed to remove file")
	ErrFileLoad   = errors.New("failed to load file")
	ErrFileSave   = errors.New("failed to save file")

	ErrGithubRequestNew       = errors.New("failed to create new request")
	ErrGithubRequestDo        = errors.New("failed to execute request")
	ErrGithubRequestStatus    = errors.New("github api returned status")
	ErrGithubRequestBodyRead  = errors.New("reading body failed")
	ErrGithubRequestBodyClose = errors.New("failed to close body")

	ErrCacheRateLimit = errors.New("failed to get ratelimit")
	ErrCacheLatest    = errors.New("failed to get latest cache")

	ErrRepoEmpty       = errors.New("repo is empty")
	ErrRepoBodyEmpty   = errors.New("body is empty")
	ErrRepoAssetsEmpty = errors.New("assets are empty")
)

// -- utils --

func NewLogFile(root *os.Root, path string) (*os.File, error) {
	timestamp := time.Now().UTC().Unix()
	filename := fmt.Sprintf("log-%d.txt", timestamp)

	file, err := root.Create(filepath.Join(path, filename))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLogFileInvalid, err)
	}

	return file, nil
}
