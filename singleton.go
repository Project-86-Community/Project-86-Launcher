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
	"context"
	"os"
	"p86l/internal/data"
	"p86l/internal/debug"
	"p86l/internal/file"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/quasilyte/gdata/v2"
)

type debugMode struct {
	IsRelease bool
	LogFile   *os.File
	Logs      bool
}

var (
	ctx = context.Background()

	TheDebugMode debugMode
	GDataM       *gdata.Manager
	l            *i18n.Locale

	e  *debug.Debug
	fs *file.AppFS
	d  *data.Data

	AppErr *debug.Error
)
