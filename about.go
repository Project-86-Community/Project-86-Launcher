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
	"p86l/assets"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type About struct {
	guigui.DefaultWidget

	background basicwidget.Background
	leadImg    basicwidget.Image
	devImg     basicwidget.Image
	leadText   basicwidget.Text
	devText    basicwidget.Text
	infoText   basicwidget.Text
}

func (a *About) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetOpacity(&a.background, 0.7)
	appender.AppendChildWidgetWithBounds(&a.background, context.Bounds(a))

	img, dErr := assets.TheImageCache.Get(e, "lead")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	a.leadImg.SetImage(img)
	img, dErr = assets.TheImageCache.Get(e, "dev")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	a.devImg.SetImage(img)

	a.infoText.SetAutoWrap(true)
	a.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	a.infoText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	a.leadText.SetScale(1.2)
	a.devText.SetScale(1.2)
	a.leadText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	a.devText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	a.leadText.SetValue(T("about.lead"))
	a.devText.SetValue(T("about.dev"))
	a.infoText.SetValue(T("about.info"))

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
		gl := layout.GridLayout{
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
		appender.AppendChildWidgetWithBounds(&a.leadText, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&a.leadImg, gl.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&a.devText, gl.CellBounds(0, 1))
		appender.AppendChildWidgetWithBounds(&a.devImg, gl.CellBounds(1, 1))
	}
	appender.AppendChildWidgetWithBounds(&a.infoText, gl.CellBounds(0, 1))

	return nil
}
