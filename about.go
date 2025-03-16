/*
 * SPDX-License-Identifier: GPL-3.0-only
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

package eightysix

import (
	"eightysix/content"
	"eightysix/internal"
	"image"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type About struct {
	guigui.DefaultWidget

	vLayout   internal.VerticalLayout
	infoText  basicwidget.Text
	leadText  basicwidget.Text
	leadImage basicwidget.Image
	devText   basicwidget.Text
	devImage  basicwidget.Image

	err error
}

func (a *About) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	img, err := content.TheImageCache.Get("lead", context.ColorMode())
	if err != nil {
		a.err = err
		return
	}
	a.leadImage.SetImage(img)
	img, err = content.TheImageCache.Get("dev", context.ColorMode())
	if err != nil {
		a.err = err
		return
	}
	a.devImage.SetImage(img)

	a.leadImage.SetSize(context, 64, 64)
	a.devImage.SetSize(context, 64, 64)

	u := float64(basicwidget.UnitSize(context))
	w, _ := a.Size(context)
	pt := guigui.Position(a).Add(image.Pt(int(0.5*u), int(0.5*u)))

	a.infoText.SetText("Welcome to Project 86 - a fan game in its early stages, \n with the primary goal of delivering a functional beta swiftly. \n We invite players to actively participate and provide feedback, \n steering the game in the right direction.")
	a.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	a.infoText.SetMultiline(true)

	a.leadText.SetText("Lead Developer - Taliayaya - Ilan Mayeux")

	a.devText.SetText("Launcher Developer - realskyquest - Sky")

	a.vLayout.SetHorizontalAlign(internal.HorizontalAlignCenter)

	a.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&a.vLayout, pt)

	a.vLayout.SetItems([]*internal.LayoutItem{
		{Widget: &a.infoText},
		{Widget: &a.leadImage},
		{Widget: &a.leadText},
		{Widget: &a.devImage},
		{Widget: &a.devText},
	})
	appender.AppendChildWidget(&a.vLayout)
}

func (a *About) Update(context *guigui.Context) error {
	if a.err != nil {
		return a.err
	}
	return nil
}

func (a *About) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(a).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
