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
	"fmt"
	"os"
	"p86l/assets/lang"
	"p86l/internal/data"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/guigui"
	"github.com/invopop/ctxi18n"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() *debug.Error {
	e = &debug.Debug{}
	fs = &file.AppFS{GdataM: GDataM}
	d = &data.Data{GDataM: GDataM}

	if TheDebugMode.IsRelease {
		logDir, err := fs.LogDir(e)
		if err.Err != nil {
			return err
		}

		if _err := os.MkdirAll(logDir, 0755); _err != nil {
			return e.New(_err, debug.FSError, debug.ErrNewDirFailed)
		}

		timestamp := time.Now().Unix()
		logFileName := fmt.Sprintf("log_%d.log", timestamp)
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, _err := os.Create(logFilePath)
		if _err != nil {
			return e.New(_err, debug.FSError, debug.ErrNewFileFailed)
		}

		TheDebugMode.LogFile = logFile

		multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	}

	d.ColorMode = guigui.ColorModeLight
	d.AppScale = 2

	if err := lang.GetLangs(); err != nil {
		return e.New(err, debug.FSError, debug.ErrFileNotFound)
	}
	ctx, err := ctxi18n.WithLocale(ctx, "en")
	if err != nil {
		panic(err)
	}
	l = ctxi18n.Locale(ctx)

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}
