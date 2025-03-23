/*
 * SPDX-License-Identifier: GPL-3.0-only
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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

package eightysix

import (
	"eightysix/internal/app"
	"eightysix/internal/content"
	"eightysix/internal/intwidget"
	"image"
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Settings struct {
	guigui.DefaultWidget

	vLayout      intwidget.VerticalLayout
	darkModeForm intwidget.Form

	darkModeText         basicwidget.Text
	darkModeToggle       basicwidget.ToggleButton
	appScaleText         basicwidget.Text
	appScaleDropdownList basicwidget.DropdownList
	openFolderButton     basicwidget.TextButton
	repairButton         basicwidget.TextButton
	clearCacheButton     basicwidget.TextButton
	deleteFilesButton    basicwidget.TextButton

	initOnce sync.Once
	err      error
}

func (s *Settings) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	s.appScaleDropdownList.SetItemsByStrings([]string{"50%", "75%", "100%", "125%", "150%"})
	s.initOnce.Do(func() {
		if content.ColorMode == guigui.ColorModeDark {
			s.darkModeToggle.SetValue(true)
		}

		s.appScaleDropdownList.SetSelectedItemIndex(content.AppScale)
	})

	s.darkModeToggle.SetOnValueChanged(func(value bool) {
		if value {
			content.ColorMode = guigui.ColorModeDark
		} else {
			content.ColorMode = guigui.ColorModeLight
		}
	})

	s.appScaleDropdownList.SetOnValueChanged(func(value int) {
		content.AppScale = value
	})

	s.openFolderButton.SetOnDown(func() {
		if content.Mgdata.ObjectPropExists("cache", "darkmode.data") {
			go func() {
				err := app.OpenFileManager(app.LauncherPathFolder())
				s.err = err
			}()
		}
	})

	s.clearCacheButton.SetOnDown(func() {
		if content.Mgdata.ObjectPropExists("cache", "darkmode.data") {
			content.Mgdata.DeleteObject("cache")
			s.darkModeToggle.SetValue(false)
			s.appScaleDropdownList.SetSelectedItemIndex(2)
		}
	})

	u := float64(basicwidget.UnitSize(context))
	w, _ := s.Size(context)
	pt := guigui.Position(s).Add(image.Pt(int(0.5*u), int(0.5*u)))

	s.darkModeText.SetText("Dark Mode")
	s.appScaleText.SetText("App Scale")
	s.openFolderButton.SetText("Open folder")
	s.repairButton.SetText("Repair")
	s.clearCacheButton.SetText("Clear cache")
	s.deleteFilesButton.SetText("Delete all files")

	s.darkModeForm.SetItems([]*intwidget.FormItem{
		{PrimaryWidget: &s.darkModeText, SecondaryWidget: &s.darkModeToggle},
	})

	s.vLayout.SetHorizontalAlign(intwidget.HorizontalAlignCenter)
	s.vLayout.SetBackground(true)
	s.vLayout.SetLineBreak(false)
	s.vLayout.SetBorder(true)

	s.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&s.vLayout, pt)

	s.vLayout.SetItems([]*intwidget.LayoutItem{
		{Widget: &s.darkModeForm},
		{Widget: &s.appScaleText},
		{Widget: &s.appScaleDropdownList},
		{Widget: &s.openFolderButton},
		{Widget: &s.repairButton},
		{Widget: &s.clearCacheButton},
		{Widget: &s.deleteFilesButton},
	})
	appender.AppendChildWidget(&s.vLayout)
}

func (s *Settings) Update(context *guigui.Context) error {
	if s.err != nil {
		return s.err
	}
	return nil
}

func (s *Settings) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(s).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
