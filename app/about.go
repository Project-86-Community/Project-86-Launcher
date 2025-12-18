/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type About struct {
	guigui.DefaultWidget

	textPanel                  basicwidget.Panel
	form                       basicwidget.Form
	text1, text2, text3, text4 basicwidget.Text
	image1, image2             aboutIcon
}

func (a *About) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&a.text1)
	adder.AddChild(&a.form)
	adder.AddChild(&a.textPanel)

	a.text1.SetValue(p86l.T("about.content"))
	a.text1.SetAutoWrap(true)

	a.text2.SetValue("Tali")
	a.text2.SetScale(1.2)

	a.text3.SetValue("Sky")
	a.text3.SetScale(1.2)

	a.text4.SetScale(0.8)
	a.text4.SetAutoWrap(true)
	a.text4.SetMultiline(true)
	a.text4.SetValue(p86l.T("about.license"))

	a.image1.setIcon(assets.LeadDeveloper)
	a.image2.setIcon(assets.DevDeveloper)

	a.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &a.text2,
			SecondaryWidget: &a.image1,
		},
		{
			PrimaryWidget:   &a.text3,
			SecondaryWidget: &a.image2,
		},
	})

	a.textPanel.SetContent(&a.text4)
	a.textPanel.SetAutoBorder(true)
	a.textPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	return nil
}

func (a *About) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &a.text1,
			},
			{
				Widget: &a.form,
			},
			{
				Widget: &a.textPanel,
				Size:   guigui.FlexibleSize(1),
			},
		},
		Gap: u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

type aboutIcon struct {
	guigui.DefaultWidget

	image basicwidget.Image

	ebitenImage *ebiten.Image
}

func (a *aboutIcon) setIcon(icon *ebiten.Image) {
	a.ebitenImage = icon
}

func (a *aboutIcon) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&a.image)
	a.image.SetImage(a.ebitenImage)

	return nil
}

func (a *aboutIcon) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &a.image,
				Size:   guigui.FixedSize(u * 2),
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (a *aboutIcon) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(u*2, u*2)
}
