/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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
	"p86l"
	"p86l/assets"
	"p86l/configs"
	"path/filepath"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"golang.org/x/text/language"
)

type Settings struct {
	guigui.DefaultWidget

	formPanel                                                                                         basicwidget.Panel
	form                                                                                              basicwidget.Form
	languageText, translateChangelogText, darkModeText, scaleText, rememberWindowText, disableBgmText basicwidget.Text
	translateChangelogToggle, darkModeToggle, rememberWindowToggle, disableBgmToggle                  basicwidget.Toggle
	languageSelect                                                                                    basicwidget.Select[language.Tag]
	scaleSegmentedControl                                                                             basicwidget.SegmentedControl[float64]
	companyText, launcherText, logsText                                                               basicwidget.Text
	companyButton, launcherButton, logsButton                                                         basicwidget.Button
	resetDataText, resetCacheText                                                                     basicwidget.Text
	resetDataButton, resetCacheButton                                                                 basicwidget.Button
}

func (s *Settings) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&s.formPanel)

	model := context.Model(s, modelKeyModel).(*p86l.Model)
	data := model.Data()
	dataFile := data.Get()
	cacheFile := model.Cache().Get()

	s.languageText.SetValue(p86l.T("settings.language"))
	s.translateChangelogText.SetValue(p86l.T("settings.translate"))
	s.darkModeText.SetValue(p86l.T("settings.darkmode"))
	s.scaleText.SetValue(p86l.T("settings.scale"))
	s.rememberWindowText.SetValue(p86l.T("settings.remember"))
	s.disableBgmText.SetValue(p86l.T("settings.backgroundm"))

	s.languageSelect.SetItems([]basicwidget.SelectItem[language.Tag]{
		{
			Text:  "English",
			Value: language.English,
		},
		{
			Text:  "French",
			Value: language.French,
		},
	})
	s.languageSelect.SetOnItemSelected(func(index int) {
		item, ok := s.languageSelect.ItemByIndex(index)
		if !ok {
			context.SetAppLocales(nil)
			return
		}

		context.SetAppLocales([]language.Tag{item.Value})
		assets.LoadLanguage(item.Value.String())
		data.Update(func(df *p86l.DataFile) {
			df.Lang = item.Value.String()
		})

		if cacheFile.Releases != nil && dataFile.TranslateChangelog && item.Value != language.English {
			model.Translate(p86l.ReleasesChangelogText(cacheFile, dataFile.UsePreRelease), item.Value.String())
		}
	})
	if !s.languageSelect.IsPopupOpen() {
		if locales := context.AppendAppLocales(nil); len(locales) > 0 {
			s.languageSelect.SelectItemByValue(locales[0])
		} else {
			s.languageSelect.SelectItemByValue(language.English)
		}
	}

	context.SetEnabled(&s.translateChangelogToggle, dataFile.Lang != "en")
	s.translateChangelogToggle.SetOnValueChanged(func(value bool) {
		data.Update(func(df *p86l.DataFile) {
			df.TranslateChangelog = value
		})

		if cacheFile.Releases != nil && value && dataFile.Lang != "en" {
			model.Translate(p86l.ReleasesChangelogText(cacheFile, dataFile.UsePreRelease), dataFile.Lang)
		}
	})
	if dataFile.TranslateChangelog {
		s.translateChangelogToggle.SetValue(true)
	} else {
		s.translateChangelogToggle.SetValue(false)
	}

	s.darkModeToggle.SetOnValueChanged(func(value bool) {
		if value {
			context.SetColorMode(guigui.ColorModeDark)
		} else {
			context.SetColorMode(guigui.ColorModeLight)
		}
		data.Update(func(df *p86l.DataFile) {
			df.UseDarkmode = value
		})
	})
	if dataFile.UseDarkmode {
		s.darkModeToggle.SetValue(true)
	} else {
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
		data.Update(func(df *p86l.DataFile) {
			df.AppScale = item.Value
		})
	})
	s.scaleSegmentedControl.SelectItemByValue(context.AppScale())

	s.rememberWindowToggle.SetOnValueChanged(func(value bool) {
		data.Update(func(df *p86l.DataFile) {
			df.Remember.Active = value
		})
	})
	if dataFile.Remember.Active {
		s.rememberWindowToggle.SetValue(true)
	} else {
		s.rememberWindowToggle.SetValue(false)
	}

	s.disableBgmToggle.SetOnValueChanged(func(value bool) {
		if value {
			model.BGMPlayer().Pause()
		} else {
			model.BGMPlayer().Play()
		}
		data.Update(func(df *p86l.DataFile) {
			df.DisableBgMusic = value
		})
	})
	if dataFile.DisableBgMusic {
		s.disableBgmToggle.SetValue(true)
	} else {
		s.disableBgmToggle.SetValue(false)
	}

	launcherPath := configs.AppName
	logsPath := filepath.Join(launcherPath, configs.FolderLogs)

	s.companyButton.SetOnDown(func() {
		model.OpenPath("")
	})
	s.launcherButton.SetOnDown(func() {
		model.OpenPath(launcherPath)
	})
	s.logsButton.SetOnDown(func() {
		model.OpenPath(logsPath)
	})

	s.companyText.SetValue(p86l.T("settings.openp86"))
	s.launcherText.SetValue(p86l.T("settings.openl"))
	s.logsText.SetValue(p86l.T("settings.openlog"))
	s.companyButton.SetText(p86l.T("common.open"))
	s.launcherButton.SetText(p86l.T("common.open"))
	s.logsButton.SetText(p86l.T("common.open"))

	s.resetDataButton.SetOnDown(model.ResetDataAsync)
	s.resetDataText.SetValue(p86l.T("settings.resetd"))
	s.resetDataButton.SetText(p86l.T("common.reset"))

	s.resetCacheButton.SetOnDown(model.ResetCacheAsync)
	s.resetCacheText.SetValue(p86l.T("settings.resetc"))
	s.resetCacheButton.SetText(p86l.T("common.reset"))

	s.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &s.languageText,
			SecondaryWidget: &s.languageSelect,
		},
		{
			PrimaryWidget:   &s.translateChangelogText,
			SecondaryWidget: &s.translateChangelogToggle,
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
		{
			PrimaryWidget:   &s.resetDataText,
			SecondaryWidget: &s.resetDataButton,
		},
		{
			PrimaryWidget:   &s.resetCacheText,
			SecondaryWidget: &s.resetCacheButton,
		},
	})

	s.formPanel.SetContent(&s.form)
	s.formPanel.SetAutoBorder(true)
	s.formPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	return nil
}

func (s *Settings) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &s.formPanel,
				Size:   guigui.FlexibleSize(1),
			},
		},
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
