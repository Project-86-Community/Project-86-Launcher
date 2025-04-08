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
	"fmt"
	"os"
	ESApp "p86l/internal/app"
	"p86l/internal/cache"
	"p86l/internal/data"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/hajimehoshi/guigui"
	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type debugMode struct {
	IsRelease bool
	LogFile   *os.File
	Logs      bool
}

var (
	TheDebugMode debugMode
	GDataM       *gdata.Manager

	AppErr        *debug.Error
	app           *ESApp.App
	githubClient  = github.NewClient(nil)
	githubContext = context.Background()
)

func Run() *debug.Error {
	app = &ESApp.App{
		Debug: &debug.Debug{},
		FS:    &file.AppFS{GdataM: GDataM},
		Data:  &data.Data{GDataM: GDataM},
		Cache: &cache.Cache{GDataM: GDataM},
	}

	if TheDebugMode.IsRelease {
		logDir, err := app.FS.LogDir(app.Debug)
		if err.Err != nil {
			return err
		}

		if _err := os.MkdirAll(logDir, 0755); _err != nil {
			return app.Debug.New(_err, debug.FSError, debug.ErrNewDirFailed)
		}

		timestamp := time.Now().Unix()
		logFileName := fmt.Sprintf("log_%d.log", timestamp)
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, _err := os.Create(logFilePath)
		if _err != nil {
			return app.Debug.New(_err, debug.FSError, debug.ErrNewFileFailed)
		}

		TheDebugMode.LogFile = logFile

		multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	}

	app.Data.ColorMode = guigui.ColorModeLight
	app.Data.AppScale = 2

	go func() {
		app.UpdateInternet()
		if app.IsInternet() {
			err := app.Cache.InitChangelog(app.Debug, githubClient, githubContext)
			if err.Err != nil {
				app.Debug.SetToast(err)
			}
		}
	}()

	return app.Debug.New(nil, debug.UnknownError, debug.ErrUnknown)
}
