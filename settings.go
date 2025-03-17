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
	"eightysix/content"
	"eightysix/internal"
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Settings struct {
	guigui.DefaultWidget

	vLayout          internal.VerticalLayout
	hLayout          internal.HorizontalLayout
	cacheButton      basicwidget.TextButton
	toggleButtonText basicwidget.Text
	toggleButton     basicwidget.ToggleButton
	repairButton     basicwidget.TextButton
	openButton       basicwidget.TextButton
	deleteButton     basicwidget.TextButton

	err error
}

func (s *Settings) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if content.Mgdata.ObjectPropExists("darkmode", "darkmode.data") {
		darkModeByte, err := content.Mgdata.LoadObjectProp("darkmode", "darkmode.data")
		if err != nil {
			s.err = err
			return
		}
		darkModeData, err := strconv.Atoi(string(darkModeByte))
		if err != nil {
			s.err = err
			return
		}
		if darkModeData == 1 {
			s.toggleButton.SetValue(true)
		}
	}

	u := float64(basicwidget.UnitSize(context))
	w, _ := s.Size(context)
	pt := guigui.Position(s).Add(image.Pt(int(0.5*u), int(0.5*u)))

	s.cacheButton.SetText("Reset cache")
	s.toggleButtonText.SetText("Dark mode")
	s.repairButton.SetText("Repair")
	s.openButton.SetText("Open folder")
	s.deleteButton.SetText("Delete all files")

	s.hLayout.SetVerticalAlign(internal.VerticalAlignStart)
	s.hLayout.DisableBackground(true)
	s.hLayout.DisableLineBreak(true)
	s.hLayout.DisableBorder(true)

	s.hLayout.SetHeight(context, 25)
	s.hLayout.SetItems([]*internal.LayoutItem{
		{Widget: &s.toggleButtonText},
		{Widget: &s.toggleButton},
	})

	s.vLayout.SetHorizontalAlign(internal.HorizontalAlignCenter)
	s.vLayout.DisableBackground(true)

	s.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&s.vLayout, pt)

	s.vLayout.SetItems([]*internal.LayoutItem{
		{Widget: &s.hLayout},
		{Widget: &s.cacheButton},
		{Widget: &s.repairButton},
		{Widget: &s.openButton},
		{Widget: &s.deleteButton},
	})
	appender.AppendChildWidget(&s.vLayout)

	s.cacheButton.SetOnDown(func() {
		content.Mgdata.DeleteObject("game")
		content.Mgdata.DeleteObject("changelog")
	})
	s.toggleButton.SetOnValueChanged(func(value bool) {
		if value {
			context.SetColorMode(guigui.ColorModeDark)
		} else {
			context.SetColorMode(guigui.ColorModeLight)
		}
		darkModeData := context.ColorMode()
		if err := content.Mgdata.SaveObjectProp("darkmode", "darkmode.data", []byte(fmt.Sprintf("%v", darkModeData))); err != nil {
			s.err = err
			return
		}
	})
	s.openButton.SetOnDown(func() {
		if content.Mgdata.ObjectPropExists("darkmode", "darkmode.data") {
			folderPath := content.Mgdata.ObjectPropPath("darkmode", "darkmode.data")
			folderPath = internal.TrimDarkModePath(folderPath)

			go func() {
				err := internal.OpenFileManager(folderPath)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
	})
	s.deleteButton.SetOnDown(func() {
		if content.Mgdata.ObjectPropExists("darkmode", "darkmode.data") {
			folderPath := content.Mgdata.ObjectPropPath("darkmode", "darkmode.data")
			folderPath = internal.TrimDarkModePath(folderPath)

			_, err := os.Stat(folderPath + "run")
			if err == nil || !os.IsNotExist(err) {
				err = os.RemoveAll(folderPath + "run")
				if err != nil {
					s.err = err
					return
				}
			}
			_, err = os.Stat(folderPath + "game.zip")
			if err == nil || !os.IsNotExist(err) {
				err = os.Remove(folderPath + "game.zip")
				if err != nil {
					s.err = err
					return
				}
			}

			s.toggleButton.SetValue(false)
			content.Mgdata.DeleteObject("darkmode")
			content.Mgdata.DeleteObject("game")
			content.Mgdata.DeleteObject("changelog")
		}
	})
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
