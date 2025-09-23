/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game Community-Game Community-Game for managing game files.
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
	"p86l"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Home struct {
	guigui.DefaultWidget

	background                                               basicwidget.Background
	form                                                     basicwidget.Form
	welcomeText, installedText, playTimeText, lastPlayedText basicwidget.Text
	usernameText, versionText, timeText, lastText            basicwidget.Text
}

func (h *Home) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&h.background)
	adder.AddChild(&h.form)

	h.welcomeText.SetValue(p86l.T("home.welcome"))
	h.usernameText.SetValue(p86l.GetUsername())
	h.installedText.SetValue(p86l.T("home.version"))
	h.versionText.SetValue("v1.8.2-alpha")
	h.playTimeText.SetValue(p86l.T("home.time"))
	h.lastPlayedText.SetValue(p86l.T("home.last"))

	h.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &h.welcomeText,
			SecondaryWidget: &h.usernameText,
		},
		{
			PrimaryWidget:   &h.installedText,
			SecondaryWidget: &h.versionText,
		},
		{
			PrimaryWidget:   &h.playTimeText,
			SecondaryWidget: &h.timeText,
		},
		{
			PrimaryWidget:   &h.lastPlayedText,
			SecondaryWidget: &h.lastText,
		},
	})

	return nil
}

func (h *Home) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Size: guigui.FlexibleSize(1),
			},
			{
				Widget: &h.background,
				Size:   guigui.FixedSize(h.form.Measure(context, guigui.FixedWidthConstraints(widgetBounds.Bounds().Dx())).Y + u),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &h.form,
						},
					},
					Padding: guigui.Padding{
						Start:  u / 2,
						Top:    u / 2,
						End:    u / 2,
						Bottom: u / 2,
					},
				},
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
