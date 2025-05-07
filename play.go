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

	background    basicwidget.Background
	actionButton  basicwidget.TextButton
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

	img, dErr := assets.TheImageCache.Get(e, "ie")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	p.websiteButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "github")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	p.githubButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "discord")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	p.discordButton.SetIcon(img)

	img, dErr = assets.TheImageCache.Get(e, "patreon")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	p.patreonButton.SetIcon(img)

	// if p.model.isInternet && p.model.cache.repo != nil {
	// 	if !p.model.play.downloading {
	// 		if dErr := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, "game.zip")); dErr.Err != nil {
	// 			context.SetEnabled(&p.actionButton, true)
	// 			p.model.play.SetStatus("install")
	// 		} else {
	// 			context.SetEnabled(&p.actionButton, true)
	// 			p.model.play.SetStatus("play")
	// 		}
	// 		p.actionButton.SetText(T(fmt.Sprintf("play.%s", p.model.play.status)))
	// 	} else {
	// 		context.SetEnabled(&p.actionButton, false)
	// 		p.actionButton.SetText(p.model.play.downloadMsg)
	// 	}
	// } else {
	// 	if dErr := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, "game.zip")); dErr.Err != nil {
	// 		context.SetEnabled(&p.actionButton, false)
	// 		p.actionButton.SetText(T("play.nointernet"))
	// 	} else {
	// 		context.SetEnabled(&p.actionButton, true)
	// 		p.actionButton.SetText(T("play.play"))
	// 	}
	// }

	p.websiteButton.SetText(T("play.website"))
	p.githubButton.SetText(T("play.github"))
	p.discordButton.SetText(T("play.discord"))
	p.patreonButton.SetText(T("play.patreon"))

	p.actionButton.SetOnDown(func() {
		// if !p.model.play.downloading {
		// 	switch p.model.play.status {
		// 	case "install":
		// 		for _, asset := range p.model.cache.repo.Assets {
		// 			if IsValidGameFile(asset.GetName()) {
		// 				p.model.play.SetDownloading(true)
		// 				go func() {
		// 					dErr := DownloadFile(p.model, "https://github.com/Taliayaya/Project-86/releases/download/v0.0.0-alpha/Project86-v0.0.0-alpha.zip", filepath.Join(fs.CompanyDirPath, "game.zip"))
		// 					if dErr != nil {
		// 						log.Error().Stack().Int("Code", dErr.Code).Str("Type", string(dErr.Type)).Err((dErr.Err)).Msg("App crashed")
		// 					}
		// 					p.model.play.SetDownloading(false)
		// 				}()
		// 			}
		// 		}
		// 	case "update":
		//
		// 	case "play":
		//
		// 	}
		// }
	})

	p.Buttons()

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(u * 4),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	{
		gl := layout.GridLayout{
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
		appender.AppendChildWidgetWithBounds(&p.websiteButton, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&p.githubButton, gl.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&p.discordButton, gl.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&p.patreonButton, gl.CellBounds(1, 1))
	}
	{
		pt := gl.Bounds.Min
		s := p.actionButton.DefaultSize(context)
		pt.X += (gl.Bounds.Dx() - s.X) / 2
		pt.Y += (gl.Bounds.Dy() - s.Y) / 2
		appender.AppendChildWidgetWithBounds(&p.actionButton, image.Rectangle{
			Min: pt.Add(image.Pt(-u*2, -u/2)),
			Max: pt.Add(s.Add(image.Pt(u*2, u/2))),
		})
	}

	return nil
}
