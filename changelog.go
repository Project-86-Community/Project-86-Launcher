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

	infoText  basicwidget.Text
	form      basicwidget.Form
	gtlText   basicwidget.Text
	gtlToggle basicwidget.Toggle
	urlButton basicwidget.TextButton

	model *Model
}

func (c *Changelog) IsChangelog() bool {
	value := c.model.cache.File()
	dErr := value.Validate(e)
	if value != nil && dErr == nil {
		return true
	}

	return false
}

func (c *Changelog) IsTranslated() bool {
	if c.model.cache.IsTranslate() && c.model.cache.TranslatedChangelog() != "" {
		return true
	}

	return false
}

func (c *Changelog) SetModel(model *Model) {
	c.model = model
}

func (c *Changelog) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if c.IsChangelog() {
		if c.IsTranslated() {
			c.infoText.SetValue(c.model.cache.TranslatedChangelog())
		} else {
			c.infoText.SetValue(c.model.cache.File().Repo.GetBody())
		}

		context.SetEnabled(&c.urlButton, true)
	} else {
		context.SetEnabled(&c.urlButton, false)
		c.infoText.SetValue("")
	}

	if c.model.networkState.InternetAvailable() {
		context.SetEnabled(&c.urlButton, true)
		if c.model.data.File().Locale == language.English.String() {
			context.SetEnabled(&c.gtlToggle, false)
			c.gtlToggle.SetValue(false)
		} else {
			context.SetEnabled(&c.gtlToggle, true)
		}

		if !c.model.cache.IsTranslate() {
			c.gtlToggle.SetValue(false)
		}
	} else {
		context.SetEnabled(&c.urlButton, false)
		context.SetEnabled(&c.gtlToggle, false)
		c.gtlToggle.SetValue(false)
	}

	c.gtlText.SetValue(T("changelog.gtl"))
	c.gtlToggle.SetOnValueChanged(func(value bool) {
		if value {
			go func() {
				result, err := t.Translate(c.model.cache.File().Repo.GetBody(), "auto", c.model.data.File().Locale)
				if err != nil {
					log.Error().Err(err).Msg("SetChangelog")
					return
				}
				c.model.cache.translatedChangelog = result.Text
				log.Info().Any("translate", result).Msg("changelog.gtlToggle")
			}()
		}
		c.model.cache.isTranslate = value
	})

	c.infoText.SetAutoWrap(true)
	c.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	c.infoText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	c.urlButton.SetText(T("changelog.open"))
	c.urlButton.SetOnDown(func() {
		if c.model.networkState.InternetAvailable() && c.IsChangelog() {
			go OpenBrowser(c.model.cache.File().Repo.GetHTMLURL())
		}
	})

	c.form.SetItems([]basicwidget.FormItem{
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
	appender.AppendChildWidgetWithBounds(&c.infoText, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&c.form, gl.CellBounds(0, 1))

	return nil
}
