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
	"p86l/assets"
	"p86l/configs"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Play struct {
	guigui.DefaultWidget

	installButton, updateButton, playButton                   basicwidget.Button
	background                                                basicwidget.Background
	form                                                      basicwidget.Form
	prereleaseText                                            basicwidget.Text
	prereleaseToggle                                          basicwidget.Toggle
	changelogText                                             basicwidget.Text
	websiteButton, githubButton, discordButton, patreonButton basicwidget.Button

	mainLayout   layout.GridLayout
	buttonLayout layout.GridLayout
	socialLayout layout.GridLayout
}

func (p *Play) Overflow(context *guigui.Context) image.Point {
	return p86l.MergeRectangles(p.mainLayout.CellBounds(0, 0), p.mainLayout.CellBounds(0, 1), p.mainLayout.CellBounds(0, 2), p.mainLayout.CellBounds(0, 3)).Size().Add(image.Pt(0, basicwidget.UnitSize(context)))
}

func (p *Play) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&p.installButton)
	appender.AppendChildWidget(&p.updateButton)
	appender.AppendChildWidget(&p.playButton)
	appender.AppendChildWidget(&p.background)
	appender.AppendChildWidget(&p.form)
	appender.AppendChildWidget(&p.websiteButton)
	appender.AppendChildWidget(&p.githubButton)
	appender.AppendChildWidget(&p.discordButton)
	appender.AppendChildWidget(&p.patreonButton)
}

func (p *Play) Build(context *guigui.Context) error {
	model := context.Model(p, modelKeyModel).(*p86l.Model)

	p.installButton.SetText("Install")
	p.updateButton.SetText("Update")
	p.playButton.SetText("Play")

	p.changelogText.SetValue(`TEST

	VER v15.15.15

	GOOD`)
	p.changelogText.SetAutoWrap(true)
	p.changelogText.SetMultiline(true)

	p.prereleaseText.SetValue("Enable Pre-release")
	p.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &p.prereleaseText,
			SecondaryWidget: &p.prereleaseToggle,
		},
		{
			PrimaryWidget: &p.changelogText,
		},
	})

	p.websiteButton.SetOnDown(func() { model.File().Open(configs.Website) })
	p.githubButton.SetOnDown(func() { model.File().Open(configs.Github) })
	p.discordButton.SetOnDown(func() { model.File().Open(configs.Discord) })
	p.patreonButton.SetOnDown(func() { model.File().Open(configs.Patreon) })

	p.websiteButton.SetIcon(assets.IE)
	p.githubButton.SetIcon(assets.Github)
	p.discordButton.SetIcon(assets.Discord)
	p.patreonButton.SetIcon(assets.Patreon)

	u := basicwidget.UnitSize(context)
	p.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u * 2),
			layout.FixedSize(u * 2),
			layout.FixedSize(p.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(p).Dx()-u)).Y + u/2),
			layout.FixedSize(int(float64(u) * 1.5)),
		},
		RowGap: u / 2,
	}
	p.buttonLayout = layout.GridLayout{
		Bounds: p.mainLayout.CellBounds(0, 1),
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(u * 4),
			layout.FixedSize(u * 4),
			layout.FixedSize(u * 4),
			layout.FlexibleSize(1),
		},
		Heights: []layout.Size{
			layout.FixedSize(u * 2),
		},
		ColumnGap: u / 2,
	}
	p.socialLayout = layout.GridLayout{
		Bounds: p.mainLayout.CellBounds(0, 3),
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(int(float64(u) * 1.5)),
			layout.FixedSize(int(float64(u) * 1.5)),
			layout.FixedSize(int(float64(u) * 1.5)),
			layout.FixedSize(int(float64(u) * 1.5)),
			layout.FlexibleSize(1),
		},
		ColumnGap: u / 2,
	}

	return nil
}

func (p *Play) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &p.installButton:
		return p.buttonLayout.CellBounds(1, 0)
	case &p.updateButton:
		return p.buttonLayout.CellBounds(2, 0)
	case &p.playButton:
		return p.buttonLayout.CellBounds(3, 0)
	case &p.background:
		return p.mainLayout.CellBounds(0, 2)
	case &p.form:
		return p.mainLayout.CellBounds(0, 2).Inset(basicwidget.UnitSize(context) / 4)
	case &p.websiteButton:
		return p.socialLayout.CellBounds(1, 0)
	case &p.githubButton:
		return p.socialLayout.CellBounds(2, 0)
	case &p.discordButton:
		return p.socialLayout.CellBounds(3, 0)
	case &p.patreonButton:
		return p.socialLayout.CellBounds(4, 0)
	}

	return image.Rectangle{}
}
