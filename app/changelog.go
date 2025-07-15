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
	"p86l"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Changelog struct {
	guigui.DefaultWidget

	form       basicwidget.Form
	text       basicwidget.Text
	viewText   basicwidget.Text
	viewButton basicwidget.Button

	model *p86l.Model
}

func (c *Changelog) SetModel(model *p86l.Model) {
	c.model = model
}

func (c *Changelog) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	c.text.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	c.text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	c.text.SetAutoWrap(true)

	if cache := c.model.Cache(); cache.IsValid() {
		if locale := c.model.Data().File().Locale; locale != "en" {
			if cache.TranslatedBody == "" {
				c.text.SetValue("...")
			} else {
				c.text.SetValue(cache.TranslatedBody)
			}
		} else {
			c.text.SetValue(cache.File().Repo.GetBody())
		}
		c.viewText.SetValue(cache.File().Repo.GetHTMLURL())
		context.SetEnabled(&c.viewButton, true)
	} else {
		c.text.SetValue("...")
		c.viewText.SetValue("?")
		context.SetEnabled(&c.viewButton, false)
	}

	c.viewButton.SetOnDown(func() {
		if value := c.viewText.Value(); value != "?" {
			go p86l.OpenBrowser(value)
		}
	})
	c.viewButton.SetText(p86l.T("changelog.view"))

	c.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget: &c.viewText,
		},
		{
			SecondaryWidget: &c.viewButton,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(c).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(c.form.DefaultSize(context).Y),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&c.text, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&c.form, gl.CellBounds(0, 1))

	return nil
}
