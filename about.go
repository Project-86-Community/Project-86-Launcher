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
	"p86l/internal/widget"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type About struct {
	guigui.DefaultWidget

	vLayout   widget.VerticalLayout
	infoText  basicwidget.Text
	leadText  basicwidget.Text
	leadImage basicwidget.Image
	devText   basicwidget.Text
	devImage  basicwidget.Image

	err *debug.Error
}

func (a *About) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	img, err := assets.TheImageCache.Get("lead")
	if err != nil {
		a.err = app.Debug.New(err, debug.FSError, debug.ErrFileNotFound)
		return
	}
	a.leadImage.SetImage(img)
	img, err = assets.TheImageCache.Get("dev")
	if err != nil {
		a.err = app.Debug.New(err, debug.FSError, debug.ErrFileNotFound)
		return
	}
	a.devImage.SetImage(img)

	a.leadImage.SetSize(context, 64, 64)
	a.devImage.SetSize(context, 64, 64)

	u := float64(basicwidget.UnitSize(context))
	w, _ := a.Size(context)
	pt := guigui.Position(a).Add(image.Pt(int(0.5*u), int(0.5*u)))

	a.infoText.SetText(WrapText(context, "Welcome to Project 86 - a fan game in its early stages, with the primary goal of delivering a functional beta swiftly. We invite players to actively participate and provide feedback, steering the game in the right direction.", w-int(1*u)))
	a.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	a.infoText.SetMultiline(true)

	a.leadText.SetText(WrapText(context, "Lead Developer - Taliayaya - Ilan Mayeux", w-int(1*u)))

	a.devText.SetText(WrapText(context, "Launcher Developer - realskyquest - Sky", w-int(1*u)))

	a.vLayout.SetHorizontalAlign(widget.HorizontalAlignCenter)

	a.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&a.vLayout, pt)

	a.vLayout.SetItems([]*widget.LayoutItem{
		{Widget: &a.infoText},
		{Widget: &a.leadImage},
		{Widget: &a.leadText},
		{Widget: &a.devImage},
		{Widget: &a.devText},
	})
	appender.AppendChildWidget(&a.vLayout)
}

func (a *About) Update(context *guigui.Context) error {
	if a.err != nil && a.err.Err != nil {
		AppErr = a.err
		return a.err.Err
	}
	return nil
}

func (a *About) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(a).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
