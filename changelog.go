/*
 * SPDX-License-Identifier: GPL-3.0-only
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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

package eightysix

import (
	"github.com/hajimehoshi/guigui"
)

type Changelog struct {
	guigui.DefaultWidget

	//vLayout         intwidget.VerticalLayout
	//changelogText   basicwidget.Text
	//vButtonLayout   intwidget.VerticalLayout
	//changelogButton basicwidget.TextButton

	//err error
}

func (c *Changelog) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	// c.changelogButton.SetOnDown(func() {
	// 	go browser.OpenURL(content.Changelog.URL)
	// })
	//
	// u := float64(basicwidget.UnitSize(context))
	// w, _ := c.Size(context)
	// pt := guigui.Position(c).Add(image.Pt(int(0.5*u), int(0.5*u)))
	//
	// c.vLayout.SetBackground(true)
	// c.vLayout.SetBorder(true)
	//
	// c.vLayout.SetWidth(context, w-int(1*u))
	// guigui.SetPosition(&c.vLayout, pt)
	//
	// changelogTextData := app.WrapText(context, content.Changelog.Body, w-int(1*u))
	// c.changelogText.SetText(changelogTextData)
	//
	// c.changelogButton.SetText("View changelog")
	// c.vButtonLayout.SetWidth(context, w-int(1*u))
	// c.vButtonLayout.SetHorizontalAlign(intwidget.HorizontalAlignCenter)
	//
	// c.vButtonLayout.SetItems([]*intwidget.LayoutItem{
	// 	{Widget: &c.changelogButton},
	// })
	//
	// c.vLayout.SetItems([]*intwidget.LayoutItem{
	// 	{Widget: &c.changelogText},
	// 	{Widget: &c.vButtonLayout},
	// })
	// appender.AppendChildWidget(&c.vLayout)
}

func (c *Changelog) Update(context *guigui.Context) error {
	// if c.err != nil {
	// 	return c.err
	// }
	return nil
}

func (c *Changelog) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(c).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
