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
	"image"
	"p86l"
	"p86l/configs"
	"path/filepath"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	background basicwidget.Background

	form1                                                     basicwidget.Form
	languageText, darkModeText, scaleText, rememberWindowText basicwidget.Text
	languageDropdownList                                      basicwidget.DropdownList[language.Tag]
	darkModeToggle                                            basicwidget.Toggle
	scaleSegmentedControl                                     basicwidget.SegmentedControl[float64]
	rememberWindowToggle                                      basicwidget.Toggle

	form2                                     basicwidget.Form
	companyText, launcherText, logsText       basicwidget.Text
	companyButton, launcherButton, logsButton basicwidget.Button

	form3            basicwidget.Form
	resetCacheText   basicwidget.Text
	resetCacheButton basicwidget.Button

	mainLayout layout.GridLayout
}

func (s *Settings) Overflow(context *guigui.Context) image.Point {
	return p86l.MergeRectangles(s.mainLayout.CellBounds(0, 0), s.mainLayout.CellBounds(0, 1), s.mainLayout.CellBounds(0, 2)).Size().Add(image.Pt(0, basicwidget.UnitSize(context)))
}

func (s *Settings) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.background)
	appender.AppendChildWidget(&s.form1)
	appender.AppendChildWidget(&s.form2)
	appender.AppendChildWidget(&s.form3)
}

func (s *Settings) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.languageText.SetValue("Language")
	s.darkModeText.SetValue("Use darkmode")
	s.scaleText.SetValue("Scale")
	s.rememberWindowText.SetValue("Remember window size & position & page")

	s.languageDropdownList.SetItems([]basicwidget.DropdownListItem[language.Tag]{
		{
			Text:  "English",
			Value: language.English,
		},
		{
			Text:  "French",
			Value: language.French,
		},
	})
	s.languageDropdownList.SetOnItemSelected(func(index int) {
		item, ok := s.languageDropdownList.ItemByIndex(index)
		if !ok {
			context.SetAppLocales(nil)
			return
		}
		if item.Value == language.English {
			context.SetAppLocales(nil)
			return
		}
		context.SetAppLocales([]language.Tag{item.Value})
	})
	if !s.languageDropdownList.IsPopupOpen() {
		if locales := context.AppendAppLocales(nil); len(locales) > 0 {
			s.languageDropdownList.SelectItemByValue(locales[0])
		} else {
			s.languageDropdownList.SelectItemByValue(language.English)
		}
	}

	s.darkModeToggle.SetOnValueChanged(func(value bool) {
		if value {
			context.SetColorMode(guigui.ColorModeDark)
		} else {
			context.SetColorMode(guigui.ColorModeLight)
		}
	})
	switch context.ColorMode() {
	case guigui.ColorModeLight:
		s.darkModeToggle.SetValue(false)
	case guigui.ColorModeDark:
		s.darkModeToggle.SetValue(true)
	default:
		s.darkModeToggle.SetValue(false)
	}

	s.scaleSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[float64]{
		{
			Text:  "50%",
			Value: 0.5,
		},
		{
			Text:  "75%",
			Value: 0.75,
		},
		{
			Text:  "100%",
			Value: 1.0,
		},
		{
			Text:  "125%",
			Value: 1.25,
		},
		{
			Text:  "150%",
			Value: 1.50,
		},
	})
	s.scaleSegmentedControl.SetOnItemSelected(func(index int) {
		item, ok := s.scaleSegmentedControl.ItemByIndex(index)
		if !ok {
			context.SetAppScale(1)
			return
		}
		context.SetAppScale(item.Value)
	})
	s.scaleSegmentedControl.SelectItemByValue(context.AppScale())

	items1 := []basicwidget.FormItem{
		{
			PrimaryWidget:   &s.languageText,
			SecondaryWidget: &s.languageDropdownList,
		},
		{
			PrimaryWidget:   &s.darkModeText,
			SecondaryWidget: &s.darkModeToggle,
		},
		{
			PrimaryWidget:   &s.scaleText,
			SecondaryWidget: &s.scaleSegmentedControl,
		},
		{
			PrimaryWidget:   &s.rememberWindowText,
			SecondaryWidget: &s.rememberWindowToggle,
		},
	}

	s.companyButton.SetOnDown(func() { model.File().Open(model.File().FS().Path()) })
	s.launcherButton.SetOnDown(func() { model.File().Open(filepath.Join(model.File().FS().Path(), configs.AppName)) })
	s.logsButton.SetOnDown(func() {
		model.File().Open(filepath.Join(model.File().FS().Path(), configs.AppName, configs.FolderLogs))
	})
	s.companyText.SetValue("Open 86-Project folder")
	s.launcherText.SetValue("Open launcher folder")
	s.logsText.SetValue("Open logs folder")
	s.companyButton.SetText("Open")
	s.launcherButton.SetText("Open")
	s.logsButton.SetText("Open")

	items2 := []basicwidget.FormItem{
		{
			PrimaryWidget:   &s.companyText,
			SecondaryWidget: &s.companyButton,
		},
		{
			PrimaryWidget:   &s.launcherText,
			SecondaryWidget: &s.launcherButton,
		},
		{
			PrimaryWidget:   &s.logsText,
			SecondaryWidget: &s.logsButton,
		},
	}

	s.resetCacheButton.SetOnDown(func() { go model.Cache().ForceRefresh() })
	s.resetCacheText.SetValue("Reset cache")
	s.resetCacheButton.SetText("Reset")

	items3 := []basicwidget.FormItem{
		{
			PrimaryWidget:   &s.resetCacheText,
			SecondaryWidget: &s.resetCacheButton,
		},
	}

	s.form1.SetItems(items1)
	s.form2.SetItems(items2)
	s.form3.SetItems(items3)
	u := basicwidget.UnitSize(context)
	s.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(s.form1.Measure(context, guigui.FixedWidthConstraints(context.Bounds(s).Dx()-u)).Y + u/2),
			layout.FixedSize(s.form2.Measure(context, guigui.FixedWidthConstraints(context.Bounds(s).Dx()-u)).Y + u/2),
			layout.FixedSize(s.form3.Measure(context, guigui.FixedWidthConstraints(context.Bounds(s).Dx()-u)).Y + u/2),
		},
	}

	return nil
}

func (s *Settings) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	switch widget {
	case &s.background:
		r1 := s.mainLayout.CellBounds(0, 0)
		r2 := s.mainLayout.CellBounds(0, 1)
		r3 := s.mainLayout.CellBounds(0, 2)
		return image.Rectangle{
			Min: r1.Min,
			Max: image.Pt(max(r1.Max.X, r2.Max.X, r3.Max.X), max(r1.Max.Y, r2.Max.Y, r3.Max.Y)),
		}
	case &s.form1:
		return s.mainLayout.CellBounds(0, 0).Inset(u / 4)
	case &s.form2:
		return s.mainLayout.CellBounds(0, 1).Inset(u / 4)
	case &s.form3:
		return s.mainLayout.CellBounds(0, 2).Inset(u / 4)
	}

	return image.Rectangle{}
}
