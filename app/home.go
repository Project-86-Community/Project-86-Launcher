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
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Home struct {
	guigui.DefaultWidget

	form                                         basicwidget.Form
	background                                   basicwidget.Background
	welcomeText, downloadedText, gameVersionText basicwidget.Text
}

func (h *Home) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&h.background)
	appender.AppendChildWidget(&h.form)
}

func (h *Home) Build(context *guigui.Context) error {
	h.welcomeText.SetValue("Welcome Test")
	h.downloadedText.SetValue("5000 Downloads")
	h.gameVersionText.SetValue("I have v15.15.15")

	h.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget: &h.welcomeText,
		},
		{
			PrimaryWidget: &h.downloadedText,
		},
		{
			PrimaryWidget: &h.gameVersionText,
		},
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(h),
		Heights: []layout.Size{
			layout.FlexibleSize(2),
			layout.FlexibleSize(1),
		},
	}
	context.SetBounds(&h.background, gl.CellBounds(0, 1), h)
	context.SetBounds(&h.form, gl.CellBounds(0, 1), h)

	return nil
}
