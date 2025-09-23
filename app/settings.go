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
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	background basicwidget.Background

	form1                                                                     basicwidget.Form
	languageText, darkModeText, scaleText, rememberWindowText, disableBgmText basicwidget.Text
	darkModeToggle, rememberWindowToggle, disableBgmToggle                    basicwidget.Toggle
	languageDropdownList                                                      basicwidget.DropdownList[language.Tag]
	scaleSegmentedControl                                                     basicwidget.SegmentedControl[float64]

	form2                                     basicwidget.Form
	companyText, launcherText, logsText       basicwidget.Text
	companyButton, launcherButton, logsButton basicwidget.Button

	form3            basicwidget.Form
	resetCacheText   basicwidget.Text
	resetCacheButton basicwidget.Button
}

func (s *Settings) Overflow(context *guigui.Context) image.Point {
	r1 := context.Bounds(&s.form1).Bounds()
	r2 := context.Bounds(&s.form2).Bounds()
	r3 := context.Bounds(&s.form3).Bounds()
	size := p86l.MergeRectangles(r1, r2, r3)

	return size.Size().Add(image.Pt(0, basicwidget.UnitSize(context)))
}

func (s *Settings) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&s.background)
	adder.AddChild(&s.form1)
	adder.AddChild(&s.form2)
	adder.AddChild(&s.form3)
}

func (s *Settings) Update(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.languageText.SetValue("Language")
	s.darkModeText.SetValue("Use darkmode")
	s.scaleText.SetValue("Scale")
	s.rememberWindowText.SetValue("Remember window size & position & page")
	s.disableBgmText.SetValue("Disable background music")

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
		model.Data().SetLang(item.Value)
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
		model.Data().SetUseDarkmode(value)
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
		model.Data().SetAppScale(item.Value)
		context.SetAppScale(item.Value)
	})
	s.scaleSegmentedControl.SelectItemByValue(context.AppScale())

	s.disableBgmToggle.SetOnValueChanged(func(value bool) {
		model.Data().SetDisableBgMusic(model.Player(), value)
	})
	switch model.Data().DisableBgMusic() {
	case true:
		s.disableBgmToggle.SetValue(true)
	case false:
		s.disableBgmToggle.SetValue(false)
	}

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
		{
			PrimaryWidget:   &s.disableBgmText,
			SecondaryWidget: &s.disableBgmToggle,
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

	return nil
}

func (s *Settings) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	return (guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &s.background,
				Size:   guigui.FixedSize(s.Overflow(context).Y - u/2),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &s.form1,
						},
						{
							Widget: &s.form2,
						},
						{
							Widget: &s.form3,
						},
					},
					Gap: u / 2,
					Padding: guigui.Padding{
						Start:  u / 4,
						Top:    u / 4,
						End:    u / 4,
						Bottom: u / 4,
					},
				},
			},
		},
	}).WidgetBounds(context, context.Bounds(s).Inset(u/4), widget)
}
