//go:generate goversioninfo -icon=../../assets/p86l.ico -manifest=app.manifest

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
	"flag"
	"fmt"
	"net"
	"os"
	"p86l"
	"p86l/app"
	"p86l/configs"
	"p86l/internal/file"
	"p86l/internal/log"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog"
)

var VERSION = "dev"

func setupLogger(fsSB *os.Root) (*zerolog.Logger, *os.File) {
	switch VERSION {
	case "dev":
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		logger := zerolog.New(output).With().Timestamp().Logger()
		for _, token := range strings.Split(os.Getenv("P86L_DEBUG"), ",") {
			if token != "log" {
				zerolog.SetGlobalLevel(zerolog.Disabled)
			}
		}
		return &logger, nil
	default:
		logFile, err := log.NewLogFile(fsSB, filepath.Join("Project-86-Launcher", "logs"))
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		multiWriter := zerolog.MultiLevelWriter(os.Stdout, logFile)
		logger := zerolog.New(multiWriter).With().Timestamp().Logger()
		return &logger, logFile
	}
}

func main() {
	port := flag.Int("instance", 54321, "Port to use for single-instance locking")
	flag.Parse()

	fs, err := file.NewFilesystem()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger, logFile := setupLogger(fs.Root())

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		logger.Error().Err(fmt.Errorf("another instance is already running (or port %d is in use): %w", *port, err)).Msg(log.NetworkManager.String())
		os.Exit(1)
	}

	model := &p86l.Model{}
	model.SetListener(listener)
	model.Log().SetLogger(logger)
	model.Log().SetLogFile(logFile)
	model.File().SetFS(fs)
	model.File().SetLogger(logger)
	model.Cache().SetLogger(logger)
	model.SetMode("home")

	model.Cache().Start()
	app := &app.Root{}
	app.SetModel(model)

	defer func() {
		if err := app.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	op := &guigui.RunOptions{
		Title:         configs.AppTitle,
		WindowMinSize: configs.AppWindowMinSize,
	}
	if err := guigui.Run(app, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
