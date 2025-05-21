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

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Play struct {
	guigui.DefaultWidget

	actionButton  basicwidget.Button
	updateButton  basicwidget.Button
	websiteButton basicwidget.Button
	githubButton  basicwidget.Button
	discordButton basicwidget.Button
	patreonButton basicwidget.Button

	model *Model
}

func (p *Play) SetModel(model *Model) {
	p.model = model
}

func (p *Play) shouldShowUpdateButton() bool {
	return p.model.lunch.status == LunchStatusUpdate
}

func (p *Play) buildImages() error {
	img, dErr := assets.TheImageCache.Get(e, "ie")
	if dErr != nil {
		gErr = dErr
		return dErr.Err
	}
	p.websiteButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "github")
	if dErr != nil {
		gErr = dErr
		return dErr.Err
	}
	p.githubButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "discord")
	if dErr != nil {
		gErr = dErr
		return dErr.Err
	}
	p.discordButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "patreon")
	if dErr != nil {
		gErr = dErr
		return dErr.Err
	}
	p.patreonButton.SetIcon(img)

	return nil
}

func (p *Play) buildButtons() {
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
	if err := p.buildImages(); err != nil {
		return err
	}

	if p.model.networkState.InternetAvailable() {
		context.SetEnabled(&p.websiteButton, true)
		context.SetEnabled(&p.githubButton, true)
		context.SetEnabled(&p.discordButton, true)
		context.SetEnabled(&p.patreonButton, true)
	} else {
		context.SetEnabled(&p.websiteButton, false)
		context.SetEnabled(&p.githubButton, false)
		context.SetEnabled(&p.discordButton, false)
		context.SetEnabled(&p.patreonButton, false)
	}

	switch p.model.lunch.status {
	case LunchStatusInstall:
		if p.model.networkState.InternetAvailable() {
			p.actionButton.SetText(T("play.install"))
		} else {
			p.actionButton.SetText(T("play.nointernet"))
		}
	case LunchStatusUpdate:
		p.actionButton.SetText(T("play.play"))
		p.updateButton.SetText(T("play.update"))
	case LunchStatusPlay:
		p.actionButton.SetText(T("play.play"))
	}

	p.websiteButton.SetText(T("play.website"))
	p.githubButton.SetText(T("play.github"))
	p.discordButton.SetText(T("play.discord"))
	p.patreonButton.SetText(T("play.patreon"))

	p.actionButton.SetOnDown(func() {

	})
	p.updateButton.SetOnDown(func() {

	})

	p.buildButtons()

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u * 4),
			layout.FlexibleSize(1),
			layout.FixedSize(u * 2),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}

	// Social buttons grid
	{
		glB := layout.GridLayout{
			Bounds: gl.CellBounds(0, 0),
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
		}
		appender.AppendChildWidgetWithBounds(&p.websiteButton, glB.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&p.githubButton, glB.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&p.discordButton, glB.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&p.patreonButton, glB.CellBounds(1, 1))
	}

	actionButtonsRow := gl.CellBounds(0, 2)
	buttonWidth := p.actionButton.DefaultSize(context).X + (u * 4)
	totalWidth := buttonWidth

	if p.shouldShowUpdateButton() {
		totalWidth += u + p.updateButton.DefaultSize(context).X
	}

	startX := actionButtonsRow.Min.X + (actionButtonsRow.Dx()-totalWidth)/2

	actionBtnBounds := image.Rect(
		startX,
		actionButtonsRow.Min.Y,
		startX+buttonWidth,
		actionButtonsRow.Max.Y,
	)
	appender.AppendChildWidgetWithBounds(&p.actionButton, actionBtnBounds)

	if p.shouldShowUpdateButton() {
		updateBtnBounds := image.Rect(
			startX+buttonWidth+u,
			actionButtonsRow.Min.Y,
			startX+totalWidth,
			actionButtonsRow.Max.Y,
		)
		appender.AppendChildWidgetWithBounds(&p.updateButton, updateBtnBounds)
	}

	return nil
}
