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
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	form                 basicwidget.Form
	localeText           basicwidget.Text
	localeDropdownList   basicwidget.DropdownList[language.Tag]
	colorModeText        basicwidget.Text
	colorModeToggle      basicwidget.Toggle
	appScaleText         basicwidget.Text
	appScaleDropdownList basicwidget.DropdownList[int]
	openFolderText       basicwidget.Text
	openFolderButton     basicwidget.TextButton
	clearDataButton      basicwidget.TextButton
	clearCacheButton     basicwidget.TextButton
	resetButton          basicwidget.TextButton

	sync  sync.Once
	model *Model
	dErr  *debug.Error
}

func (s *Settings) SetModel(model *Model) {
	s.model = model
}

func (s *Settings) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.localeText.SetValue(T("settings.locale"))
	s.colorModeText.SetValue(T("settings.colormode"))
	s.appScaleText.SetValue(T("settings.appscale"))

	s.openFolderText.SetValue(T("settings.openfoldertext"))
	s.openFolderButton.SetText(T("settings.openfolder"))

	s.clearDataButton.SetText(T("settings.resetdata"))
	s.clearCacheButton.SetText(T("settings.resetcache"))
	s.resetButton.SetText(T("settings.reset"))

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
		data := s.model.data.File()
		item, ok := s.localeDropdownList.ItemByIndex(index)
		if !ok {
			data.Locale = language.English.String()
			s.dErr = s.model.data.SetData(context, data)
			return
		}
		data.Locale = item.ID.String()
		s.dErr = s.model.data.SetData(context, data)
		s.model.cache.SetIsTranslate(false)
		s.model.cache.SetTranslatedChangelog("")
	})

	s.colorModeToggle.SetOnValueChanged(func(value bool) {
		data := s.model.data.File()
		if value {
			data.ColorMode = guigui.ColorModeDark
			s.dErr = s.model.data.SetData(context, data)
		} else {
			data.ColorMode = guigui.ColorModeLight
			s.dErr = s.model.data.SetData(context, data)
		}
	})

	s.appScaleDropdownList.SetItems([]basicwidget.DropdownListItem[int]{
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
	s.appScaleDropdownList.SetOnItemSelected(func(index int) {
		data := s.model.data.File()
		item, ok := s.appScaleDropdownList.ItemByIndex(index)
		if !ok {
			data.AppScale = 2
			s.dErr = s.model.data.SetData(context, data)
			return
		}
		data.AppScale = item.ID
		s.dErr = s.model.data.SetData(context, data)
	})

	s.openFolderButton.SetOnDown(func() {
		if dErr := fs.OpenFileManager(e, filepath.Join(fs.CompanyDirPath, configs.AppName)); dErr != nil {
			e.SetToast(dErr)
		}
	})

	s.clearDataButton.SetOnDown(func() {
		s.colorModeToggle.SetValue(false)
		s.localeDropdownList.SelectItemByID(language.English)
		s.appScaleDropdownList.SelectItemByID(2)
	})
	s.clearCacheButton.SetOnDown(func() {
		//s.model.cache.repo = nil
	})
	s.resetButton.SetOnDown(func() {

	})

	s.sync.Do(func() {
		if context.ColorMode() == guigui.ColorModeDark {
			s.colorModeToggle.SetValue(true)
		} else {
			s.colorModeToggle.SetValue(false)
		}
		data := s.model.data.File()
		locale, err := language.Parse(data.Locale)
		if err != nil {
			s.dErr = e.New(err, debug.DataError, debug.ErrDataLoad)
		}
		s.localeDropdownList.SelectItemByID(locale)
		s.appScaleDropdownList.SelectItemByID(s.model.data.GetAppScale(context.AppScale()))
	})

	if s.dErr != nil {
		gErr = s.dErr
		return s.dErr.Err
	}

	s.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &s.localeText,
			SecondaryWidget: &s.localeDropdownList,
		},
		{
			PrimaryWidget:   &s.colorModeText,
			SecondaryWidget: &s.colorModeToggle,
		},
		{
			PrimaryWidget:   &s.appScaleText,
			SecondaryWidget: &s.appScaleDropdownList,
		},
		{
			PrimaryWidget:   &s.openFolderText,
			SecondaryWidget: &s.openFolderButton,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &s.clearDataButton,
		},
		{
			PrimaryWidget:   nil,
			SecondaryWidget: &s.clearCacheButton,
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
					return layout.FixedSize(0)
				}
				return layout.FixedSize(s.form.DefaultSize(context).Y)
			}),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.form, gl.CellBounds(0, 0))

	return nil
}
