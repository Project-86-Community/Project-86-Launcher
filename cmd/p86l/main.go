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
	"fmt"
	"os"
	"p86l"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/debug"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var AppBuild string

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if AppBuild == "release" {
		p86l.TheDebugMode.IsRelease = true
		p86l.TheDebugMode.Logs = true
	} else {
		for _, token := range strings.Split(os.Getenv("P86L_DEBUG"), ",") {
			switch token {
			case "logs":
				log.Logger = log.Output(zerolog.ConsoleWriter{
					Out:        os.Stderr,
					TimeFormat: "2006/01/02 15:04:05",
				})
				p86l.TheDebugMode.IsRelease = false
				p86l.TheDebugMode.Logs = true
			}
		}
	}

	if !p86l.TheDebugMode.Logs {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

func main() {
	appName := fmt.Sprintf("%s/%s", configs.CompanyName, configs.AppName)
	if runtime.GOOS == "windows" {
		appName = fmt.Sprintf("%s\\%s", configs.CompanyName, configs.AppName)
	}

	m, _err := gdata.Open(gdata.Config{
		AppName: appName,
	})
	if _err != nil {
		log.Panic().Int("Code", debug.ErrGDataOpenFailed).Str("Type", string(debug.FSError)).Err(_err).Send()
	}

	iconImages, _err := assets.GetIconImages()
	if _err != nil {
		log.Panic().Int("Code", debug.ErrIconNotFound).Str("Type", string(debug.FSError)).Err(_err).Send()
	}

	p86l.GDataM = m

	log.Info().Str("Detected OS", runtime.GOOS).Send()
	ebiten.SetWindowIcon(iconImages)

	op := &guigui.RunOptions{
		Title:           "Project 86 Launcher",
		WindowMinWidth:  500,
		WindowMinHeight: 280,
	}
	if _err = guigui.Run(&p86l.Root{}, op); _err != nil {
		log.Error().Stack().Int("Code", p86l.AppErr.Code).Str("Type", string(p86l.AppErr.Type)).Err(p86l.AppErr.Err).Msg("App crashed")
		os.Exit(1)
	}
}
