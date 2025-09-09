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

package app

import (
	"image"
	"net"
	"os"
	"p86l"
	"p86l/assets"
	"p86l/internal/file"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

type Root struct {
	guigui.DefaultWidget

	backgroundImage basicwidget.Image
	background      basicwidget.Background
	sidebar         Sidebar
	home            Home
	play            Play
	settings        Settings
	about           About

	model                   p86l.Model
	backgroundImageSize     image.Point
	backgroundImagePosition image.Point
	mainLayout              layout.GridLayout
}

func (r *Root) handleBackgroundImage(context *guigui.Context) {
	imgWidth := assets.Banner.Bounds().Dx()
	imgHeight := assets.Banner.Bounds().Dy()

	windowBounds := context.Bounds(r)
	windowWidth := windowBounds.Dx()
	windowHeight := windowBounds.Dy()

	imgAspectRatio := float64(imgWidth) / float64(imgHeight)
	windowAspectRatio := float64(windowWidth) / float64(windowHeight)

	var newWidth, newHeight int
	var xOffset, yOffset int

	if imgAspectRatio > windowAspectRatio {
		// The image is wider than the window. Scale by height and crop width.
		newHeight = windowHeight
		newWidth = int(float64(windowHeight) * imgAspectRatio)
		xOffset = (windowWidth - newWidth) / 2
		yOffset = 0
	} else {
		// The image is taller than the window. Scale by width and crop height.
		newWidth = windowWidth
		newHeight = int(float64(windowWidth) / imgAspectRatio)
		xOffset = 0
		yOffset = (windowHeight - newHeight) / 2
	}

	r.backgroundImageSize = image.Pt(newWidth, newHeight)
	r.backgroundImagePosition = image.Pt(xOffset, yOffset)
}

func (r *Root) SetListener(listener net.Listener) {
	r.model.SetListener(listener)
}

func (r *Root) SetLog(logger zerolog.Logger, logFile *os.File) {
	log := r.model.Log()
	log.SetLogger(logger)
	log.SetLogFile(logFile)
}

func (r *Root) SetFS(fs *file.Filesystem) {
	r.model.File().SetFS(fs)
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return &r.model
	default:
		return nil
	}
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.backgroundImage)
	appender.AppendChildWidget(&r.sidebar)
	if mode := r.model.Mode(); mode == "home" {
		appender.AppendChildWidget(&r.home)
	} else {
		appender.AppendChildWidget(&r.background)
		switch mode {
		case "play":
			appender.AppendChildWidget(&r.play)
		case "settings":
			appender.AppendChildWidget(&r.settings)
		case "about":
			appender.AppendChildWidget(&r.about)
		}
	}
}

func (r *Root) Build(context *guigui.Context) error {
	r.backgroundImage.SetImage(assets.Banner)
	r.handleBackgroundImage(context)
	context.SetOpacity(&r.background, 0.9)

	u := basicwidget.UnitSize(context)
	r.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * u),
			layout.FlexibleSize(1),
		},
	}

	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.backgroundImage:
		return image.Rect(
			r.backgroundImagePosition.X,
			r.backgroundImagePosition.Y,
			r.backgroundImagePosition.X+r.backgroundImageSize.X,
			r.backgroundImagePosition.Y+r.backgroundImageSize.Y,
		)
	case &r.background:
		return r.mainLayout.CellBounds(1, 0)
	case &r.sidebar:
		return r.mainLayout.CellBounds(0, 0)
	case &r.home:
		return r.mainLayout.CellBounds(1, 0)
	case &r.play:
		return r.mainLayout.CellBounds(1, 0)
	case &r.settings:
		return r.mainLayout.CellBounds(1, 0)
	case &r.about:
		return r.mainLayout.CellBounds(1, 0)
	}

	return image.Rectangle{}
}

func (r *Root) Close() error {
	return r.model.Close()
}
