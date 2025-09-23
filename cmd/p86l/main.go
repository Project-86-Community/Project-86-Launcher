//go:generate goversioninfo -icon=../../assets/images/icon.ico -manifest=app.manifest

/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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
	"flag"
	"fmt"
	"net"
	"p86l"
	"p86l/app"
	"p86l/configs"
	"p86l/internal/log"

	"github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"
)

var VERSION = "dev"

func main() {
	port := flag.Int("instance", 54321, "Port to use for single-instance locking")
	flag.Parse()

	root, model, fs, logger, logFile, err := app.NewRoot(VERSION)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = fs.Close(); _ = logFile.Close() }()

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		logger.Fatal().Err(fmt.Errorf("another instance is already running (or port %d is in use): %w", *port, err)).Msg(log.NetworkManager.String())
	}
	defer func() { _ = listener.Close() }()

	images, err := p86l.GetIcons()
	if err != nil {
		logger.Fatal().Err(err).Msg(log.ErrorManager.String())
	}
	ebiten.SetWindowIcon(images)

	op := &guigui.RunOptions{
		Title:         configs.AppTitle,
		WindowMinSize: configs.AppWindowMinSize,
	}
	if err := guigui.Run(root, op); err != nil {
		logger.Fatal().Err(err).Msg(log.ErrorManager.String())
	}
	logger.Info().Str(log.Lifecycle, "application closing, saving data...").Msg(log.AppManager.String())
	model.Stop()
	logger.Info().Str(log.Lifecycle, "application closed successfully").Msg(log.AppManager.String())
}
