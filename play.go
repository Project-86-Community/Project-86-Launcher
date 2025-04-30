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
	"p86l/configs"
	"p86l/internal/debug"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog/log"
)

type Play struct {
	guigui.DefaultWidget

	background    basicwidget.Background
	actionButton  basicwidget.TextButton
	form          basicwidget.Form
	websiteButton basicwidget.TextButton
	githubButton  basicwidget.TextButton
	discordButton basicwidget.TextButton
	patreonButton basicwidget.TextButton

	model *Model
}

func (p *Play) SetModel(model *Model) {
	p.model = model
}

func (p *Play) Buttons() {
	p.websiteButton.SetOnDown(func() {
		go OpenBrowser(configs.Website)
	})
	p.githubButton.SetOnDown(func() {
		go OpenBrowser(configs.Github)
	})
	p.discordButton.SetOnDown(func() {
		go OpenBrowser(configs.Discord)
	})
	p.patreonButton.SetOnDown(func() {
		go OpenBrowser(configs.Patreon)
	})
}

func (p *Play) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetOpacity(&p.background, 0.7)
	appender.AppendChildWidgetWithBounds(&p.background, context.Bounds(p))

	img, err := assets.TheImageCache.Get("ie")
	if err != nil {
		aErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
	}
	p.websiteButton.SetImage(img)

	img, err = assets.TheImageCache.Get("github")
	if err != nil {
		aErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
	}
	p.githubButton.SetImage(img)

	img, err = assets.TheImageCache.Get("discord")
	if err != nil {
		aErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
	}
	p.discordButton.SetImage(img)

	img, err = assets.TheImageCache.Get("patreon")
	if err != nil {
		aErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
	}
	p.patreonButton.SetImage(img)

	if p.model.isInternet == true {
		context.SetEnabled(&p.actionButton, true)
		p.actionButton.SetText("Play")
	} else {
		context.SetEnabled(&p.actionButton, false)
		p.actionButton.SetText("No connection!")
	}

	p.websiteButton.SetText("Website - Offical website!")
	p.githubButton.SetText("Github - Help develop the game?")
	p.discordButton.SetText("Discord - Meet the community!")
	p.patreonButton.SetText("Patreon - Donate to the devs?")

	p.actionButton.SetOnDown(func() {
		log.Info()
	})

	p.Buttons()

	u := basicwidget.UnitSize(context)
	for i, bounds := range (layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u * 4),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}).CellBounds() {
		switch i {
		case 0:
			for j, innerBounds := range (layout.GridLayout{
				Bounds: bounds,
				Widths: []layout.Size{
					layout.FlexibleSize(1),
					layout.FlexibleSize(1),
				},
				Heights: []layout.Size{
					layout.FixedSize(u * 2),
					layout.FixedSize(u * 2),
				},
				RowGap:    u / 2,
				ColumnGap: u / 2,
			}).CellBounds() {
				switch j {
				case 0:
					appender.AppendChildWidgetWithBounds(&p.websiteButton, innerBounds)
				case 1:
					appender.AppendChildWidgetWithBounds(&p.githubButton, innerBounds)
				case 2:
					appender.AppendChildWidgetWithBounds(&p.discordButton, innerBounds)
				case 3:
					appender.AppendChildWidgetWithBounds(&p.patreonButton, innerBounds)
				}
			}
		case 1:
			pt := bounds.Min
			s := p.actionButton.DefaultSize(context)
			pt.X += (bounds.Dx() - s.X) / 2
			pt.Y += (bounds.Dy() - s.Y) / 2
			appender.AppendChildWidgetWithBounds(&p.actionButton, image.Rectangle{
				Min: pt.Add(image.Pt(-u*2, -u/2)),
				Max: pt.Add(s.Add(image.Pt(u*2, u/2))),
			})
		}
	}

	return nil
}
