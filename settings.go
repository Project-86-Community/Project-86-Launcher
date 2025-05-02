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

	background           basicwidget.Background
	form                 basicwidget.Form
	localeText           basicwidget.Text
	localeDropdownList   basicwidget.DropdownList[language.Tag]
	colorModeText        basicwidget.Text
	colorModeToggle      basicwidget.Toggle
	appScaleText         basicwidget.Text
	appScaleDropdownList basicwidget.DropdownList[int]
	openFolderText       basicwidget.Text
	openFolderButton     basicwidget.TextButton
	clearDataText        basicwidget.Text
	clearDataButton      basicwidget.TextButton
	clearCacheText       basicwidget.Text
	clearCacheButton     basicwidget.TextButton
	resetText            basicwidget.Text
	resetButton          basicwidget.TextButton

	sync  sync.Once
	model *Model
	dErr  *debug.Error
}

func (s *Settings) SetModel(model *Model) {
	s.model = model
}

func (s *Settings) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetOpacity(&s.background, 0.7)
	appender.AppendChildWidgetWithBounds(&s.background, context.Bounds(s))

	s.localeText.SetText(l.T("settings.locale"))
	s.localeDropdownList.SetItems([]basicwidget.DropdownListItem[language.Tag]{
		{
			Text: "English",
			Tag:  language.English,
		},
		{
			Text: "French",
			Tag:  language.French,
		},
		{
			Text: "Japanese",
			Tag:  language.Japanese,
		},
		{
			Text: "Korean",
			Tag:  language.Korean,
		},
	})
	s.localeDropdownList.SetOnItemSelected(func(index int) {
		item, ok := s.localeDropdownList.ItemByIndex(index)
		if !ok {
			s.dErr = s.model.data.SetLocale(context, language.English)
			return
		}
		s.dErr = s.model.data.SetLocale(context, item.Tag)
	})

	s.colorModeText.SetText(l.T("settings.colormode"))
	s.colorModeToggle.SetOnValueChanged(func(value bool) {
		if value == true {
			s.dErr = s.model.data.SetColorMode(context, guigui.ColorModeDark)
		} else {
			s.dErr = s.model.data.SetColorMode(context, guigui.ColorModeLight)
		}
	})

	s.appScaleText.SetText(l.T("settings.appscale"))
	s.appScaleDropdownList.SetItems([]basicwidget.DropdownListItem[int]{
		{
			Text: "50%",
			Tag:  0,
		},
		{
			Text: "75%",
			Tag:  1,
		},
		{
			Text: "100%",
			Tag:  2,
		},
		{
			Text: "125%",
			Tag:  3,
		},
		{
			Text: "150%",
			Tag:  4,
		},
	})
	s.appScaleDropdownList.SetOnItemSelected(func(index int) {
		item, ok := s.appScaleDropdownList.ItemByIndex(index)
		if !ok {
			s.dErr = s.model.data.SetAppScale(context, 2)
			return
		}
		s.dErr = s.model.data.SetAppScale(context, item.Tag)
	})

	s.openFolderText.SetText(l.T("settings.openfoldertext"))
	s.openFolderButton.SetText(l.T("settings.openfolder"))
	s.openFolderButton.SetOnDown(func() {
		dErr := fs.OpenFileManager(e, filepath.Join(fs.CompanyDirPath, configs.AppName))
		e.SetToast(dErr)
	})

	s.clearDataText.SetText(l.T("settings.cleardatatext"))
	s.clearDataButton.SetText(l.T("settings.clear"))
	s.clearDataButton.SetOnDown(func() {

	})

	s.clearCacheText.SetText(l.T("settings.clearcachetext"))
	s.clearCacheButton.SetText(l.T("settings.clear"))
	s.clearCacheButton.SetOnDown(func() {

	})

	s.resetText.SetText(l.T("settings.resettext"))
	s.resetButton.SetText(l.T("settings.reset"))
	s.resetButton.SetOnDown(func() {

	})

	s.sync.Do(func() {
		s.localeDropdownList.SelectItemByTag(s.model.data.locale)
		if context.ColorMode() == guigui.ColorModeDark {
			s.colorModeToggle.SetValue(true)
		} else {
			s.colorModeToggle.SetValue(false)
		}
		s.appScaleDropdownList.SelectItemByTag(s.model.data.GetAppScale(context.AppScale()))
	})

	if s.dErr != nil && s.dErr.Err != nil {
		aErr = s.dErr
		return aErr.Err
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
			PrimaryWidget:   &s.clearDataText,
			SecondaryWidget: &s.clearDataButton,
		},
		{
			PrimaryWidget:   &s.clearCacheText,
			SecondaryWidget: &s.clearCacheButton,
		},
		{
			PrimaryWidget:   &s.resetText,
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
	for i, bounds := range gl.RepeatingCellBounds() {
		if i >= 1 {
			break
		}
		appender.AppendChildWidgetWithBounds(&s.form, bounds)
	}

	return nil
}
