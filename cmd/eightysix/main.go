/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Ilan Mayeux
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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
	"eightysix"
	"eightysix/assets"
	"eightysix/internal/content"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/quasilyte/gdata/v2"
)

var AppBuild string

func init() {
	if AppBuild == "release" {
		content.TheDebugMode.Logs = true
	} else {
		for _, token := range strings.Split(os.Getenv("P86L_DEBUG"), ",") {
			switch token {
			case "logs":
				content.TheDebugMode.Logs = true
			}
		}
	}
}

func main() {
	var appName string
	if runtime.GOOS == "windows" {
		appName = "Project-86-Community\\Project-86-Launcher"
	} else {
		appName = "Project-86-Community/Project-86-Launcher"
	}

	m, err := gdata.Open(gdata.Config{
		AppName: appName,
	})
	if err != nil {
		panic(err)
	}
	content.Mgdata = m

	iconImages, err := assets.GetIconImages()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowIcon(iconImages)
	op := &guigui.RunOptions{
		Title:           "Project 86 Launcher",
		WindowMinWidth:  500,
		WindowMinHeight: 280,
	}
	if err = guigui.Run(&eightysix.Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
