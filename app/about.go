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
)

type About struct {
	guigui.DefaultWidget

	background                 basicwidget.Background
	form                       basicwidget.Form
	text1, text2, text3, text4 basicwidget.Text
	image1, image2             aboutIcon
}

func (a *About) Overflow(context *guigui.Context) image.Point {
	r1 := context.Bounds(&a.text1).Bounds()
	r2 := context.Bounds(&a.form).Bounds()
	r3 := context.Bounds(&a.text4).Bounds()
	size := p86l.MergeRectangles(r1, r2, r3)

	return size.Size().Add(image.Pt(0, basicwidget.UnitSize(context)))
}

func (a *About) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&a.text1)
	adder.AddChild(&a.background)
	adder.AddChild(&a.form)
	adder.AddChild(&a.text4)
}

func (a *About) Update(context *guigui.Context) error {
	a.text1.SetValue("A Launcher developed for Project-86 for managing game files.")
	a.text1.SetAutoWrap(true)

	a.text2.SetValue("Tali")
	a.text2.SetScale(1.2)

	a.text3.SetValue("Sky")
	a.text3.SetScale(1.2)

	a.text4.SetValue(`Background music: "Project: 86 OST: Legion" by XYETRY.

 Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
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

	return nil
}

func (a *About) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	return (guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &a.text1,
			},
			{
				Widget: &a.background,
				Size:   guigui.FixedSize(a.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(a).Dx()-u)).Y + u/2),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &a.form,
						},
					},
					Padding: guigui.Padding{
						Start:  u / 4,
						Top:    u / 4,
						End:    u / 4,
						Bottom: u / 4,
					},
				},
			},
			{
				Widget: &a.text4,
			},
		},
		Gap: u / 2,
	}).WidgetBounds(context, context.Bounds(a).Inset(u/2), widget)
}

type aboutIcon struct {
	guigui.DefaultWidget

	image basicwidget.Image

	ebitenImage *ebiten.Image
}

func (a *aboutIcon) setIcon(icon *ebiten.Image) {
	a.ebitenImage = icon
}

func (a *aboutIcon) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&a.image)
}

func (a *aboutIcon) Update(context *guigui.Context) error {
	a.image.SetImage(a.ebitenImage)

	return nil
}

func (a *aboutIcon) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	return (guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &a.image,
				Size:   guigui.FixedSize(u * 3),
			},
		},
	}).WidgetBounds(context, context.Bounds(a), widget)
}

func (a *aboutIcon) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(u*3, u*3)
}
