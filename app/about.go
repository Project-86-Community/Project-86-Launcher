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

type About struct {
	guigui.DefaultWidget

	leadImg      basicwidget.Image
	devImg       basicwidget.Image
	leadText     basicwidget.Text
	devText      basicwidget.Text
	aboutContent aboutContent
}

func (a *About) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	img1, err1 := assets.TheImageCache.Get(p86l.E, "lead")
	img2, err2 := assets.TheImageCache.Get(p86l.E, "dev")

	if err := cmp.Or(err1, err2); err != nil {
		p86l.GErr = err
		return err.Err
	}

	a.leadImg.SetImage(img1)
	a.devImg.SetImage(img2)

	a.leadText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	a.leadText.SetScale(1.2)
	a.leadText.SetValue(p86l.T("about.lead"))

	a.devText.SetScale(1.2)
	a.devText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	a.devText.SetValue(p86l.T("about.dev"))

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(a).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	{
		glP := layout.GridLayout{
			Bounds: gl.CellBounds(0, 0),
			Widths: []layout.Size{
				layout.FlexibleSize(2),
				layout.FlexibleSize(1),
			},
			Heights: []layout.Size{
				layout.FlexibleSize(1),
				layout.FlexibleSize(1),
			},
			RowGap:    u / 2,
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&a.leadText, glP.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&a.leadImg, glP.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&a.devText, glP.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&a.devImg, glP.CellBounds(1, 1))
	}
	appender.AppendChildWidgetWithBounds(&a.aboutContent, gl.CellBounds(0, 1))

	return nil
}

type aboutContent struct {
	guigui.DefaultWidget

	infoText    basicwidget.Text
	licenseText basicwidget.Text
}

func (a *aboutContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	a.infoText.SetAutoWrap(true)
	a.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	a.infoText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	a.infoText.SetValue(p86l.T("about.info"))

	a.licenseText.SetAutoWrap(true)
	a.licenseText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	a.licenseText.SetVerticalAlign(basicwidget.VerticalAlignBottom)
	a.licenseText.SetScale(0.6)
	context.SetOpacity(&a.licenseText, 0.5)
	a.licenseText.SetValue(p86l.ALicense)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(a),
		Heights: []layout.Size{
			layout.FlexibleSize(2),
			layout.FlexibleSize(1),
			layout.FixedSize(u),
		},
	}
	appender.AppendChildWidgetWithBounds(&a.infoText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&a.licenseText, gl.CellBounds(0, 1))

	return nil
}
