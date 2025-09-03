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
