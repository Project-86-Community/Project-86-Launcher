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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type About struct {
	guigui.DefaultWidget

	background                 basicwidget.Background
	form                       basicwidget.Form
	text1, text2, text3, text4 basicwidget.Text
	image1, image2             aboutIcon

	mainLayout layout.GridLayout
}

func (a *About) Overflow(context *guigui.Context) image.Point {
	return p86l.MergeRectangles(a.mainLayout.CellBounds(0, 0), a.mainLayout.CellBounds(0, 1), a.mainLayout.CellBounds(0, 2)).Size().Add(image.Pt(0, basicwidget.UnitSize(context)))
}

func (a *About) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&a.text1)
	appender.AppendChildWidget(&a.background)
	appender.AppendChildWidget(&a.form)
	appender.AppendChildWidget(&a.text4)
}

func (a *About) Build(context *guigui.Context) error {
	a.text1.SetValue("A Launcher developed for Project-86 for managing game files.")
	a.text1.SetAutoWrap(true)

	a.text2.SetValue("Tali")
	a.text2.SetScale(1.2)

	a.text3.SetValue("Sky")
	a.text3.SetScale(1.2)

	a.text4.SetValue(`Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 Copyright (C) 2025 Project 86 Community

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.`)
	a.text4.SetScale(0.8)
	a.text4.SetAutoWrap(true)
	a.text4.SetMultiline(true)

	a.image1.setIcon(assets.LeadDeveloper)
	a.image2.setIcon(assets.DevDeveloper)

	items := []basicwidget.FormItem{
		{
			PrimaryWidget:   &a.text2,
			SecondaryWidget: &a.image1,
		},
		{
			PrimaryWidget:   &a.text3,
			SecondaryWidget: &a.image2,
		},
	}

	a.form.SetItems(items)
	u := basicwidget.UnitSize(context)
	a.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(a).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(a.text1.Measure(context, guigui.FixedWidthConstraints(context.Bounds(a).Dx()-u)).Y + u),
			layout.FixedSize(a.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(a).Dx()-u)).Y + u/2),
			layout.FixedSize(a.text4.Measure(context, guigui.FixedWidthConstraints(context.Bounds(a).Dx()-u)).Y),
		},
		RowGap: u / 2,
	}

	return nil
}

func (a *About) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &a.text1:
		return a.mainLayout.CellBounds(0, 0)
	case &a.background:
		return a.mainLayout.CellBounds(0, 1)
	case &a.form:
		return a.mainLayout.CellBounds(0, 1).Inset(basicwidget.UnitSize(context) / 4)
	case &a.text4:
		return a.mainLayout.CellBounds(0, 2)
	}

	return image.Rectangle{}
}

type aboutIcon struct {
	guigui.DefaultWidget

	image basicwidget.Image

	mainLayout  layout.GridLayout
	ebitenImage *ebiten.Image
}

func (a *aboutIcon) setIcon(icon *ebiten.Image) {
	a.ebitenImage = icon
}

func (a *aboutIcon) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&a.image)
}

func (a *aboutIcon) Build(context *guigui.Context) error {
	a.image.SetImage(a.ebitenImage)
	u := basicwidget.UnitSize(context)
	a.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(a),
		Heights: []layout.Size{
			layout.FixedSize(u * 3),
		},
	}

	return nil
}

func (a *aboutIcon) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &a.image:
		return a.mainLayout.CellBounds(0, 0)
	}

	return image.Rectangle{}
}

func (a *aboutIcon) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(u*3, u*3)
}
