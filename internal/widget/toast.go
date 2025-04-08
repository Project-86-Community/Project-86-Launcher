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

package widget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Toast struct {
	guigui.DefaultWidget

	toastPanel        basicwidget.ScrollablePanel
	toastText         basicwidget.Text
	toastCloseButton  basicwidget.TextButton
	widthMinusDefault int
}

func (t *Toast) SetOnDown(callback func()) {
	t.toastCloseButton.SetOnDown(callback)
}

func (t *Toast) SetText(text string) {
	t.toastText.SetText(text)
}

func (t *Toast) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	u := float64(basicwidget.UnitSize(context))
	w, h := t.Size(context)
	p := guigui.Position(t)

	t.toastCloseButton.SetText("Close")

	t.toastPanel.SetSize(context, w-int(6*u), h-int(0.5*u))
	t.toastPanel.SetContent(func(context *guigui.Context, childAppender *basicwidget.ContainerChildWidgetAppender, offsetX, offsetY float64) {
		p := guigui.Position(&t.toastPanel).Add(image.Pt(int(offsetX), int(offsetY)))

		guigui.SetPosition(&t.toastText, image.Pt(p.X, p.Y))
		childAppender.AppendChildWidget(&t.toastText)
	})
	guigui.SetPosition(&t.toastPanel, p.Add(image.Pt(int(0.5*u), int(0.2*u))))

	t.toastCloseButton.SetWidth(int(4 * u))
	guigui.SetPosition(&t.toastCloseButton, p.Add(image.Pt(w-int(4.5*u), int(0.3*u))))

	appender.AppendChildWidget(&t.toastPanel)
	appender.AppendChildWidget(&t.toastCloseButton)
}

func (t *Toast) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := guigui.Bounds(t)

	// Draw background
	bgColor := basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.875)
	basicwidget.DrawRoundedRect(context, dst, bounds, bgColor, 1)
}

func (t *Toast) SetWidth(context *guigui.Context, width int) {
	t.widthMinusDefault = width - defaultFormWidth(context)
}

func (t *Toast) Size(context *guigui.Context) (int, int) {
	width := t.widthMinusDefault + defaultFormWidth(context)
	return width, t.height(context)
}

func (t *Toast) height(context *guigui.Context) int {
	return int(1.5 * float64(basicwidget.UnitSize(context)))
}
