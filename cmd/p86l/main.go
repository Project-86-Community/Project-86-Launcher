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

package main

import (
	"os"
	"p86l"
	"p86l/cmd/p86l/p86li"
	"p86l/configs"
	"runtime"

	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	p86li.Run()

	if !p86l.TheDebugMode.Logs {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

func main() {
	log.Info().Str("Detected OS", runtime.GOOS).Send()

	op := &guigui.RunOptions{
		Title:         configs.AppTitle,
		WindowMinSize: configs.AppWindowMinSize,
	}
	if err := guigui.Run(&p86l.Root{}, op); err != nil {
		appErr := p86l.GetAppError()
		log.Error().Stack().Int("Code", appErr.Code).Str("Type", string(appErr.Type)).Err((appErr.Err)).Msg("App crashed")
		os.Exit(1)
	}
}
