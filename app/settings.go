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

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	background                            basicwidget.Background
	form                                  basicwidget.Form
	languageText, darkModeText, scaleText basicwidget.Text
	languageDropdownList                  basicwidget.DropdownList[language.Tag]
	darkModeToggle                        basicwidget.Toggle
	scaleSegmentedControl                 basicwidget.SegmentedControl[float64]

	mainLayout layout.GridLayout
}

func (s *Settings) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.background)
	appender.AppendChildWidget(&s.form)
}

func (s *Settings) Build(context *guigui.Context) error {
	s.languageText.SetValue("Language")
	s.darkModeText.SetValue("Use darkmode")
	s.scaleText.SetValue("Scale")

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

	items := []basicwidget.FormItem{
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
	}

	s.form.SetItems(items)
	u := basicwidget.UnitSize(context)
	s.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(s.form.Measure(context, guigui.FixedWidthConstraints(context.Bounds(s).Dx()-u)).Y + u/2),
		},
	}

	return nil
}

func (s *Settings) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &s.background:
		return s.mainLayout.CellBounds(0, 0)
	case &s.form:
		return s.mainLayout.CellBounds(0, 0).Inset(basicwidget.UnitSize(context) / 4)
	}

	return image.Rectangle{}
}
