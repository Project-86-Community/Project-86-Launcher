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
	"p86l"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Home struct {
	guigui.DefaultWidget

	background                                           basicwidget.Background
	form                                                 basicwidget.Form
	welcomeText, usernameText, downloadText, versionText basicwidget.Text

	mainLayout layout.GridLayout
}

func (h *Home) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&h.background)
	adder.AddChild(&h.form)
}

func (h *Home) Update(context *guigui.Context) error {
	h.welcomeText.SetValue("Welcome Test")
	h.usernameText.SetValue(p86l.GetUsername())
	h.downloadText.SetValue("Downloaded")
	h.versionText.SetValue("v15.15.15")

	h.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &h.welcomeText,
			SecondaryWidget: &h.usernameText,
		},
		{
			PrimaryWidget:   &h.downloadText,
			SecondaryWidget: &h.versionText,
		},
	})

	u := basicwidget.UnitSize(context)
	h.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(h).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(2),
			layout.FixedSize(h.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(h).Dx()-u)).Y + u/2),
		},
	}

	return nil
}

func (h *Home) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	switch widget {
	case &h.background:
		return h.mainLayout.CellBounds(0, 1)
	case &h.form:
		return h.mainLayout.CellBounds(0, 1).Inset(u / 4)
	}

	return image.Rectangle{}
}
