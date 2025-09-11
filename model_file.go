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

package p86l

import (
	"p86l/internal/file"
	"p86l/internal/log"
	"p86l/internal/open"

	"github.com/rs/zerolog"
)

type FileModel struct {
	logger *zerolog.Logger
	fs     *file.Filesystem
}

func (f *FileModel) FS() *file.Filesystem {
	return f.fs
}

//

func (f *FileModel) SetLogger(logger *zerolog.Logger) {
	f.logger = logger
}

func (f *FileModel) SetFS(fs *file.Filesystem) {
	f.fs = fs
}

// -- common --

func (f *FileModel) Open(input string) {
	if err := open.Open(input); err != nil {
		f.logger.Warn().Str("FileModel", "Open").Err(err).Msg(log.ErrorManager.String())
	}
}

func (f *FileModel) Close() error {
	return f.fs.Close()
}
