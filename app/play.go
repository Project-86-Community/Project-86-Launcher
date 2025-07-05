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
	"cmp"
	"p86l"
	"p86l/assets"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Play struct {
	guigui.DefaultWidget

	content playContent
	links   playLinks

	model *p86l.Model
}

func (p *Play) SetModel(model *p86l.Model) {
	p.model = model
}

func (p *Play) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(p.links.DefaultSize(context).Y),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&p.content, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&p.links, gl.CellBounds(0, 1))

	return nil
}

type playContent struct {
	guigui.DefaultWidget

	actionButton basicwidget.Button
}

func (p *playContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.actionButton.SetText(p86l.T("play.play"))
	
	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(p),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(2 * u),
			layout.FlexibleSize(1),
		},
		Widths: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&p.actionButton, gl.CellBounds(1, 1))

	return nil
}

type playLinks struct {
	guigui.DefaultWidget

	websiteButton basicwidget.Button
	githubButton  basicwidget.Button
	discordButton basicwidget.Button
	patreonButton basicwidget.Button
}

func (p *playLinks) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	img1, err1 := assets.TheImageCache.Get(p86l.E, "ie")
	img2, err2 := assets.TheImageCache.Get(p86l.E, "github")
	img3, err3 := assets.TheImageCache.Get(p86l.E, "discord")
	img4, err4 := assets.TheImageCache.Get(p86l.E, "patreon")

	if err := cmp.Or(err1, err2, err3, err4); err != nil {
		p86l.GErr = err
		return err.Err
	}

	p.websiteButton.SetIcon(img1)
	p.githubButton.SetIcon(img2)
	p.discordButton.SetIcon(img3)
	p.patreonButton.SetIcon(img4)

	p.websiteButton.SetText(p86l.T("play.website"))
	p.githubButton.SetText(p86l.T("play.github"))
	p.discordButton.SetText(p86l.T("play.discord"))
	p.patreonButton.SetText(p86l.T("play.patreon"))

	u := basicwidget.UnitSize(context)
	var gl layout.GridLayout
	if context.AppSize().X >= 1280 {
		gl = layout.GridLayout{
			Bounds: context.Bounds(p),
			Heights: []layout.Size{
				layout.FixedSize(u * 2),
				layout.FixedSize(u * 2),
				layout.FixedSize(u * 2),
				layout.FixedSize(u * 2),
			},
			Widths: []layout.Size{
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
			},
			RowGap:    u / 2,
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&p.websiteButton, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&p.githubButton, gl.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&p.discordButton, gl.CellBounds(2, 0))
		appender.AppendChildWidgetWithBounds(&p.patreonButton, gl.CellBounds(3, 0))
	} else if context.AppSize().X >= 1024 {
		gl = layout.GridLayout{
			Bounds: context.Bounds(p),
			Heights: []layout.Size{
				layout.FixedSize(u * 2),
				layout.FixedSize(u * 2),
			},
			Widths: []layout.Size{
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
			},
			RowGap:    u / 2,
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&p.websiteButton, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&p.githubButton, gl.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&p.discordButton, gl.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&p.patreonButton, gl.CellBounds(1, 1))
	} else {
		gl = layout.GridLayout{
			Bounds: context.Bounds(p),
			Heights: []layout.Size{
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
			},
			RowGap:    u / 2,
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&p.websiteButton, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&p.githubButton, gl.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&p.discordButton, gl.CellBounds(0, 2))
		appender.AppendChildWidgetWithBounds(&p.patreonButton, gl.CellBounds(0, 3))
	}

	return nil
}
