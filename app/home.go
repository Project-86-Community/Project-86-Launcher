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
	formPanel1, formPanel2                                   basicwidget.Panel
	form1, form2                                             basicwidget.Form
	welcomeText, installedText, playTimeText, lastPlayedText basicwidget.Text
	usernameText, versionText, timeText, lastText            basicwidget.Text
}

func (h *Home) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&h.background)
	adder.AddChild(&h.formPanel1)
	adder.AddChild(&h.formPanel2)

	model := context.Model(h, modelKeyModel).(*p86l.Model)
	dataFile := model.Data().Get()

	h.usernameText.SetAutoWrap(true)
	h.installedText.SetAutoWrap(true)
	h.timeText.SetAutoWrap(true)
	h.lastText.SetAutoWrap(true)

	h.welcomeText.SetValue(p86l.T("home.welcome"))
	h.usernameText.SetValue(p86l.GetUsername())
	h.installedText.SetValue(p86l.T("home.version"))
	if dataFile.UsePreRelease {
		h.versionText.SetValue(dataFile.InstalledPreRelease)
	} else {
		h.versionText.SetValue(dataFile.InstalledGame)
	}
	h.playTimeText.SetValue(p86l.T("home.time"))
	h.lastPlayedText.SetValue(p86l.T("home.last"))

	h.form1.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &h.welcomeText,
			SecondaryWidget: &h.usernameText,
		},
		{
			PrimaryWidget:   &h.installedText,
			SecondaryWidget: &h.versionText,
		},
	})

	h.form2.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &h.playTimeText,
			SecondaryWidget: &h.timeText,
		},
		{
			PrimaryWidget:   &h.lastPlayedText,
			SecondaryWidget: &h.lastText,
		},
	})

	h.formPanel1.SetContent(&h.form1)
	h.formPanel1.SetAutoBorder(true)
	h.formPanel1.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	h.formPanel2.SetContent(&h.form2)
	h.formPanel2.SetAutoBorder(true)
	h.formPanel2.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

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
				Size:   guigui.FixedSize(h.Measure(context, guigui.FixedWidthConstraints(widgetBounds.Bounds().Dx())).Y - int(float64(u)*2.5)),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionHorizontal,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &h.formPanel1,
							Size:   guigui.FlexibleSize(3),
						},
						{
							Widget: &h.formPanel2,
							Size:   guigui.FlexibleSize(2),
						},
					},
					Gap: u / 2,
					Padding: guigui.Padding{
						Start:  u / 4,
						Top:    u / 4,
						End:    u / 4,
						Bottom: u / 4,
					},
				},
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
