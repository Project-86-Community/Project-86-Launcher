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
	"p86l/internal/debug"
	"p86l/internal/widget"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/pkg/browser"
)

type Changelog struct {
	guigui.DefaultWidget

	vLayout         widget.VerticalLayout
	changelogText   basicwidget.Text
	vButtonLayout   widget.VerticalLayout
	changelogButton basicwidget.TextButton
}

func (c *Changelog) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	c.changelogButton.SetOnDown(func() {
		go func() {
			if err := browser.OpenURL(app.Cache.Changelog.URL); err != nil {
				app.Debug.SetToast(app.Debug.New(err, debug.AppError, debug.ErrBrowserOpen))
			}
		}()
	})

	u := float64(basicwidget.UnitSize(context))
	w, _ := c.Size(context)
	pt := guigui.Position(c).Add(image.Pt(int(0.5*u), int(0.5*u)))

	c.vLayout.SetBackground(true)
	c.vLayout.SetBorder(true)

	c.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&c.vLayout, pt)

	if app.Cache.Changelog != nil {
		changelogTextData := WrapText(context, app.Cache.Changelog.Body, w-int(1*u))
		c.changelogText.SetText(changelogTextData)
	} else {
		c.changelogText.SetText("NO INTERNET")
	}

	c.changelogButton.SetText("View changelog")
	c.vButtonLayout.SetWidth(context, w-int(1*u))
	c.vButtonLayout.SetHorizontalAlign(widget.HorizontalAlignCenter)

	c.vButtonLayout.SetItems([]*widget.LayoutItem{
		{Widget: &c.changelogButton},
	})

	c.vLayout.SetItems([]*widget.LayoutItem{
		{Widget: &c.changelogText},
		{Widget: &c.vButtonLayout},
	})
	appender.AppendChildWidget(&c.vLayout)
}

func (c *Changelog) Update(context *guigui.Context) error {
	return nil
}

func (c *Changelog) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(c).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
