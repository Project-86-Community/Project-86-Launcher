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
	urlButton basicwidget.Button

	model *Model
}

func (c *Changelog) SetModel(model *Model) {
	c.model = model
}

func (c *Changelog) IsChangelog() bool {
	file := c.model.CacheM.File()
	return file.Validate(e) == nil
}

func (c *Changelog) IsTranslated() bool {
	return c.model.CacheM.IsTranslate() && c.model.CacheM.TranslatedChangelog() != ""
}

func (c *Changelog) handleTranslationToggle(value bool) {
	if value {
		go c.translateChangelog()
	}
	c.model.CacheM.isTranslate = value
}

func (c *Changelog) translateChangelog() {
	body := c.model.CacheM.File().Repo.GetBody()
	targetLang := c.model.DataM.File().Locale
	result, err := t.Translate(body, "auto", targetLang)
	if err != nil {
		log.Error().Err(err).Msg("Translation failed")
		return
	}
	c.model.CacheM.translatedChangelog = result.Text
	log.Info().Any("translation", result).Msg("Changelog translated")
}

func (c *Changelog) configureInfoText() {
	c.infoText.SetAutoWrap(true)
	c.infoText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	c.infoText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
}

func (c *Changelog) configureURLButton(context *guigui.Context) {
	c.urlButton.SetText(T("changelog.open"))
	c.urlButton.SetOnDown(func() {
		if c.model.networkState.InternetAvailable() && c.IsChangelog() {
			go OpenBrowser(c.model.CacheM.File().Repo.GetHTMLURL())
		}
	})
	c.form.SetItems([]basicwidget.FormItem{
		{PrimaryWidget: &c.gtlText, SecondaryWidget: &c.gtlToggle},
		{SecondaryWidget: &c.urlButton},
	})
}

func (c *Changelog) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	isChangelog := c.IsChangelog()
	hasNet := c.model.networkState.InternetAvailable()
	locale := c.model.DataM.File().Locale

	if isChangelog {
		if c.IsTranslated() {
			c.infoText.SetValue(c.model.CacheM.TranslatedChangelog())
		} else {
			c.infoText.SetValue(c.model.CacheM.File().Repo.GetBody())
		}
	} else {
		c.infoText.SetValue("")
	}

	urlEnabled := isChangelog && hasNet
	context.SetEnabled(&c.urlButton, urlEnabled)

	gtlEnabled := hasNet && locale != language.English.String() && isChangelog
	context.SetEnabled(&c.gtlToggle, gtlEnabled)
	c.gtlToggle.SetValue(gtlEnabled && c.model.CacheM.IsTranslate())

	c.gtlText.SetValue(T("changelog.gtl"))
	c.gtlToggle.SetOnValueChanged(c.handleTranslationToggle)
	c.configureInfoText()
	c.configureURLButton(context)

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
