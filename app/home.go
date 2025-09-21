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
)

type Home struct {
	guigui.DefaultWidget

	background                                           basicwidget.Background
	form                                                 basicwidget.Form
	welcomeText, usernameText, downloadText, versionText basicwidget.Text
}

func (h *Home) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&h.background)
	adder.AddChild(&h.form)
}

func (h *Home) Update(context *guigui.Context) error {
	h.welcomeText.SetValue("Welcome Test")
	h.usernameText.SetValue(p86l.GetUsername())
	h.downloadText.SetValue("Downloaded")
	h.versionText.SetValue("v1.8.2-alpha")

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

	return nil
}

func (h *Home) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	return (guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Size: guigui.FlexibleSize(1),
			},
			{
				Widget: &h.background,
				Size:   guigui.FixedSize(h.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(h).Dx()-u)).Y + u/2),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &h.form,
						},
					},
					Padding: guigui.Padding{
						Start:  u / 4,
						Top:    u / 4,
						End:    u / 4,
						Bottom: u / 4,
					},
				},
			},
		},
	}).WidgetBounds(context, context.Bounds(h).Inset(u/2), widget)
}
