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
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	form                  basicwidget.Form
	localeText            basicwidget.Text
	localeDropdownList    basicwidget.DropdownList[language.Tag]
	colorModeText         basicwidget.Text
	colorModeToggle       basicwidget.Toggle
	scaleText             basicwidget.Text
	scaleSegmentedControl basicwidget.SegmentedControl[int]
	openFolderText        basicwidget.Text
	openFolderButton      basicwidget.TextButton
	resetDataButton       basicwidget.TextButton
	resetCacheButton      basicwidget.TextButton
	resetButton           basicwidget.TextButton

	model *Model
	dErr  *debug.Error
}

func (s *Settings) SetModel(model *Model) {
	s.model = model
}

func (s *Settings) assertErr(dErr *debug.Error) {
	if dErr != nil {
		s.dErr = dErr
	}
}

func (s *Settings) buildLocale(context *guigui.Context) {
	s.localeDropdownList.SetItems([]basicwidget.DropdownListItem[language.Tag]{
		{
			Text: "English",
			ID:   language.English,
		},
		{
			Text: "French",
			ID:   language.French,
		},
		{
			Text: "Japanese",
			ID:   language.Japanese,
		},
	})
	s.localeDropdownList.SetOnItemSelected(func(index int) {
		item, ok := s.localeDropdownList.ItemByIndex(index)
		if !ok {
			context.SetAppLocales(nil)
			return
		}
		if item.ID == language.English {
			s.assertErr(s.model.data.SetLocale(context, language.English))
			context.SetAppLocales(nil)
			return
		}
		s.assertErr(s.model.data.SetLocale(context, item.ID))
		context.SetAppLocales([]language.Tag{item.ID})
		s.model.cache.isTranslate = false
		s.model.cache.translatedChangelog = ""
	})
	if !s.localeDropdownList.IsPopupOpen() {
		if locales := context.AppendAppLocales(nil); len(locales) > 0 {
			s.localeDropdownList.SelectItemByID(locales[0])
		} else {
			s.localeDropdownList.SelectItemByID(language.English)
		}
	}
}

func (s *Settings) buildAppScale(context *guigui.Context) {
	s.scaleSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[int]{
		{
			Text: "50%",
			ID:   0,
		},
		{
			Text: "75%",
			ID:   1,
		},
		{
			Text: "100%",
			ID:   2,
		},
		{
			Text: "125%",
			ID:   3,
		},
		{
			Text: "150%",
			ID:   4,
		},
	})
	s.scaleSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := s.scaleSegmentedControl.ItemByIndex(index)
		if !ok {
			s.assertErr(s.model.data.SetAppScale(context, 2))
			return
		}
		s.assertErr(s.model.data.SetAppScale(context, item.ID))
	})
	s.scaleSegmentedControl.SelectItemByID(s.model.data.File().AppScale)
}

func (s *Settings) buildColorMode(context *guigui.Context) {
	s.colorModeToggle.SetOnValueChanged(func(value bool) {
		if value {
			s.assertErr(s.model.data.SetColorMode(context, guigui.ColorModeDark))
		} else {
			s.assertErr(s.model.data.SetColorMode(context, guigui.ColorModeLight))
		}
	})
	switch context.ColorMode() {
	case guigui.ColorModeLight:
		s.colorModeToggle.SetValue(false)
	case guigui.ColorModeDark:
		s.colorModeToggle.SetValue(true)
	default:
		s.colorModeToggle.SetValue(false)
	}
}

func (s *Settings) buildOpenFolder() {
	s.openFolderButton.SetOnDown(func() {
		go func() {
			if dErr := fs.OpenFileManager(e, filepath.Join(fs.CompanyDirPath, configs.AppName)); dErr != nil {
				e.SetToast(dErr)
			}
		}()
	})
}

func (s *Settings) buildReset() {
	s.resetDataButton.SetOnDown(func() {
		s.colorModeToggle.SetValue(false)
		s.localeDropdownList.SelectItemByID(language.English)
		s.scaleSegmentedControl.SelectItemByID(2)
	})
	s.resetCacheButton.SetOnDown(func() {
		s.model.cache.isTranslate = false
		s.model.cache.translatedChangelog = ""
		var cacheFile file.Cache
		s.model.cache.SetCache(&cacheFile)
	})
}

func (s *Settings) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.localeText.SetValue(T("settings.locale"))
	s.buildLocale(context)

	s.scaleText.SetValue(T("settings.appscale"))
	s.buildAppScale(context)

	s.colorModeText.SetValue(T("settings.colormode"))
	s.buildColorMode(context)

	s.openFolderText.SetValue(T("settings.openfoldertext"))
	s.openFolderButton.SetText(T("settings.openfolder"))
	s.buildOpenFolder()

	s.resetDataButton.SetText(T("settings.resetdata"))
	s.resetCacheButton.SetText(T("settings.resetcache"))
	s.resetButton.SetText(T("settings.reset"))
	s.buildReset()

	if s.dErr != nil {
		gErr = s.dErr
		return s.dErr.Err
	}

	s.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &s.localeText,
			SecondaryWidget: &s.localeDropdownList,
		},
		{
			PrimaryWidget:   &s.colorModeText,
			SecondaryWidget: &s.colorModeToggle,
		},
		{
			PrimaryWidget:   &s.scaleText,
			SecondaryWidget: &s.scaleSegmentedControl,
		},
		{
			PrimaryWidget:   &s.openFolderText,
			SecondaryWidget: &s.openFolderButton,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &s.resetDataButton,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &s.resetCacheButton,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &s.resetButton,
		},
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.LazySize(func(row int) layout.Size {
				if row >= 1 {
					return layout.FixedSize(1)
				}
				return layout.FixedSize(s.form.DefaultSize(context).Y)
			}),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.form, gl.CellBounds(0, 0))

	return nil
}
