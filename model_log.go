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
	"os"

	"github.com/rs/zerolog"
)

type LogModel struct {
	logger  zerolog.Logger
	logFile *os.File
}

func (l *LogModel) Logger() zerolog.Logger {
	return l.logger
}

func (l *LogModel) LogFile() *os.File {
	return l.logFile
}

func (l *LogModel) SetLogger(logger zerolog.Logger) {
	l.logger = logger
}

func (l *LogModel) SetLogFile(logFile *os.File) {
	l.logFile = logFile
}

// -- common --

func (l *LogModel) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
