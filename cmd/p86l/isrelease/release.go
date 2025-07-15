//go:build release
// +build release

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

package isrelease

import (
	"os"
	"p86l"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() {
	e := &debug.Debug{}
	a, dErr := file.NewFS(e)
	if dErr != nil {
		log.Error().Stack().Int("Code", dErr.Code).Str("Type", string(dErr.Type)).Err((dErr.Err)).Msg("Run")
	}

	logFile, err := a.Root.Create(filepath.Join(a.DirAppPath(), "log.txt"))
	if err != nil {
		log.Error().Stack().Int("Code", dErr.Code).Str("Type", string(dErr.Type)).Err((dErr.Err)).Msg("Run")
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	p86l.TheDebugMode.Logs = true
}
