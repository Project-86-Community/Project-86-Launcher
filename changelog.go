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
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Changelog struct {
	guigui.DefaultWidget

	background basicwidget.Background
	infoText   basicwidget.Text
	form       basicwidget.Form
	gtlText    basicwidget.Text
	gtlToggle  basicwidget.Toggle
	urlButton  basicwidget.TextButton

	isTranslate bool
	model       *Model
}

func (c *Changelog) SetModel(model *Model) {
	c.model = model
}

func (c *Changelog) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetOpacity(&c.background, 0.7)
	appender.AppendChildWidgetWithBounds(&c.background, context.Bounds(c))

	if c.model.cache.repo.GetBody() != "" {
		if c.isTranslate && c.model.cache.translatedChangelog != "" {
			c.infoText.SetText(c.model.cache.translatedChangelog)
		} else {
			c.infoText.SetText(c.model.cache.repo.GetBody())
		}
	} else {
		c.infoText.SetText("")
	}
	if c.model.cache.repo.GetHTMLURL() != "" {
		context.SetEnabled(&c.urlButton, true)
	} else {
		context.SetEnabled(&c.urlButton, false)
	}
	if c.model.data.locale == language.English {
		context.SetEnabled(&c.gtlToggle, false)
		c.gtlToggle.SetValue(false)
	} else {
		context.SetEnabled(&c.gtlToggle, true)
	}

	c.gtlText.SetText("Translate via Google TL")
	c.gtlToggle.SetOnValueChanged(func(value bool) {
		if value {
			go func() {
				result, err := t.Translate(c.model.cache.repo.GetBody(), "auto", c.model.data.locale.String())
				if err != nil {
					log.Error().Err(err).Msg("SetChangelog")
					return
				}
				c.model.cache.translatedChangelog = result.Text
				log.Info().Any("translate", result).Msg("changelog.gtlToggle")
			}()
		}
		c.isTranslate = value
	})

	c.infoText.SetAutoWrap(true)
	c.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	c.infoText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	c.urlButton.SetText("Open")
	c.urlButton.SetOnDown(func() {
		if c.model.cache.repo.GetBody() != "" {
			go OpenBrowser(c.model.cache.repo.GetBody())
		}
	})

	c.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &c.gtlText,
			SecondaryWidget: &c.gtlToggle,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &c.urlButton,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(c).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}
	for i, bounds := range gl.CellBounds() {
		switch i {
		case 0:
			appender.AppendChildWidgetWithBounds(&c.infoText, bounds)
		case 1:
			appender.AppendChildWidgetWithBounds(&c.form, bounds)
		}
	}

	return nil
}
