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
	"image"
	"p86l/assets"
	"p86l/internal/debug"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.RootWidget

	background basicwidget.Background
	bgImage    basicwidget.Image
	sidebar    Sidebar
	home       Home
	changelog  Changelog
	back       Back

	model Model
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	img, err := assets.TheImageCache.Get("banner")
	if err != nil {
		AppErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
	}
	r.bgImage.SetImage(img)

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	aspectRatio := float64(imgHeight) / float64(imgWidth)

	windowSize := context.Size(r)
	availableWidth := windowSize.X

	newHeight := int(float64(availableWidth) * aspectRatio)

	if newHeight < windowSize.Y {
		newHeight = windowSize.Y
		availableWidth = int(float64(newHeight) / aspectRatio)
	}

	context.SetSize(&r.bgImage, image.Pt(availableWidth+2, newHeight+2))

	yOffset := 0
	if newHeight > windowSize.Y {
		yOffset = -(newHeight - windowSize.Y) / 2
	}
	appender.AppendChildWidgetWithPosition(&r.bgImage, image.Pt(00, yOffset))

	r.sidebar.SetModel(&r.model)

	for i, bounds := range (layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}).CellBounds() {
		switch i {
		case 0:
			context.SetOpacity(&r.sidebar, 0.7)
			appender.AppendChildWidgetWithBounds(&r.sidebar, bounds)
		case 1:
			switch r.model.Mode() {
			case "home":
				appender.AppendChildWidgetWithBounds(&r.home, bounds)
			case "changelog":
				appender.AppendChildWidgetWithBounds(&r.changelog, bounds)
			default:
				appender.AppendChildWidgetWithBounds(&r.back, bounds)
			}
		}
	}

	return nil
}
