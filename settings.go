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
	"p86l/internal/debug"
	"p86l/internal/file"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog/log"
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
	openFolderButton      basicwidget.Button
	resetDataButton       basicwidget.Button
	resetCacheButton      basicwidget.Button
	resetButton           basicwidget.Button

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
			s.assertErr(s.model.DataM.SetLocale(context, language.English))
			context.SetAppLocales(nil)
			return
		}
		s.assertErr(s.model.DataM.SetLocale(context, item.ID))
		context.SetAppLocales([]language.Tag{item.ID})
		s.model.CacheM.isTranslate = false
		s.model.CacheM.translatedChangelog = ""
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
			s.assertErr(s.model.DataM.SetAppScale(context, 2))
			return
		}
		s.assertErr(s.model.DataM.SetAppScale(context, item.ID))
	})
	s.scaleSegmentedControl.SelectItemByID(s.model.DataM.File().AppScale)
}

func (s *Settings) buildColorMode(context *guigui.Context) {
	s.colorModeToggle.SetOnValueChanged(func(value bool) {
		if value {
			s.assertErr(s.model.DataM.SetColorMode(context, guigui.ColorModeDark))
		} else {
			s.assertErr(s.model.DataM.SetColorMode(context, guigui.ColorModeLight))
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
		log.Info().Str("Open Folder", fs.DirAppPath()).Msg("Settings")
		go func() {
			if dErr := fs.OpenFileManager(e, fs.DirAppPath()); dErr != nil {
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
		s.model.CacheM.isTranslate = false
		s.model.CacheM.translatedChangelog = ""
		var cacheFile file.Cache
		s.model.CacheM.SetCache(cacheFile)
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

func (s *Settings) HandleButtonInput(context *guigui.Context) guigui.HandleInputResult {
	if s.dErr != nil {
		gErr = s.dErr
		return guigui.AbortHandlingInputByWidget(s)
	}

	currentIndex := s.scaleSegmentedControl.SelectedItemIndex()
	itemsCount := 5

	if currentIndex >= 0 && currentIndex < itemsCount && context.IsFocusedOrHasFocusedChild(&s.scaleSegmentedControl) {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
			newIndex := currentIndex - 1
			if newIndex >= 0 {
				s.scaleSegmentedControl.SelectItemByIndex(newIndex)
				if item, ok := s.scaleSegmentedControl.ItemByIndex(newIndex); ok && item.ID != s.model.DataM.File().AppScale {
					s.assertErr(s.model.DataM.SetAppScale(context, item.ID))
				}
				return guigui.HandleInputByWidget(s)
			}
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
			newIndex := currentIndex + 1
			if newIndex < itemsCount {
				s.scaleSegmentedControl.SelectItemByIndex(newIndex)
				if item, ok := s.scaleSegmentedControl.ItemByIndex(newIndex); ok && item.ID != s.model.DataM.File().AppScale {
					s.assertErr(s.model.DataM.SetAppScale(context, item.ID))
				}
				return guigui.HandleInputByWidget(s)
			}
		}
	}

	return guigui.HandleInputResult{}
}

func (s *Settings) Tick(context *guigui.Context) error {
	if s.dErr != nil {
		gErr = s.dErr
		return s.dErr.Err
	}

	if context.IsWidgetHitAtCursor(&s.scaleSegmentedControl) || context.IsFocusedOrHasFocusedChild(&s.scaleSegmentedControl) {
		_, dy := ebiten.Wheel()

		currentIndex := s.scaleSegmentedControl.SelectedItemIndex()
		itemsCount := 5

		newIndex := currentIndex - int(dy)

		if newIndex < 0 {
			newIndex = 0
		} else if newIndex >= itemsCount {
			newIndex = itemsCount - 1
		}

		if newIndex != currentIndex {
			s.scaleSegmentedControl.SelectItemByIndex(newIndex)
			if item, ok := s.scaleSegmentedControl.ItemByIndex(newIndex); ok && item.ID != s.model.DataM.File().AppScale {
				s.assertErr(s.model.DataM.SetAppScale(context, item.ID))
			}
			context.SetFocused(&s.scaleSegmentedControl, true)
		}
	}

	return nil
}
