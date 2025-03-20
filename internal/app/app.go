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

package app

import (
	"eightysix"
	"eightysix/content"
	"eightysix/internal"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Root struct {
	guigui.RootWidget

	lastCheckInternet    time.Time
	checkInternetTimeout time.Duration

	sidebar  eightysix.Sidebar
	home     eightysix.Home
	settings eightysix.Settings
	//changelog eightysix.Changelog
	//about     eightysix.About

	initOnce sync.Once
	err      error
}

func (r *Root) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	r.initOnce.Do(func() {
		r.checkInternetTimeout = 2 * time.Second

		if content.Mgdata.ObjectPropExists("cache", "darkmode.data") {
			darkModeByte, err := content.Mgdata.LoadObjectProp("cache", "darkmode.data")
			if err != nil {
				r.err = err
				return
			}
			darkModeData, err := strconv.Atoi(string(darkModeByte))
			if err != nil {
				r.err = err
				return
			}
			content.ColorMode = guigui.ColorMode(darkModeData)
		} else {
			internal.SaveDarkMode()
		}

		if content.Mgdata.ObjectPropExists("cache", "appscale.data") {
			appScaleByte, err := content.Mgdata.LoadObjectProp("cache", "appscale.data")
			if err != nil {
				r.err = err
				return
			}
			appScaleData, err := strconv.Atoi(string(appScaleByte))
			if err != nil {
				r.err = err
				return
			}
			content.AppScale = appScaleData
		} else {
			internal.SaveAppScale()
		}
	})

	appender.AppendChildWidget(&r.sidebar)

	guigui.SetPosition(&r.sidebar, guigui.Position(r))
	sw, _ := r.sidebar.Size(context)
	p := guigui.Position(r)
	p.X += sw
	guigui.SetPosition(&r.home, p)
	guigui.SetPosition(&r.settings, p)
	//guigui.SetPosition(&r.changelog, p)
	//guigui.SetPosition(&r.about, p)

	switch r.sidebar.SelectedItemTag() {
	case "home":
		appender.AppendChildWidget(&r.home)
	case "settings":
		appender.AppendChildWidget(&r.settings)
	case "changelog":
		//appender.AppendChildWidget(&r.changelog)
	case "about":
		//appender.AppendChildWidget(&r.about)
	}
}

func (r *Root) Update(context *guigui.Context) error {
	if r.err != nil {
		return r.err
	}

	if content.ColorMode != context.ColorMode() {
		context.SetColorMode(content.ColorMode)
		internal.SaveDarkMode()
	}
	if content.AppScale != internal.GetAppScale(context.AppScale()) {
		internal.SetAppScale(context)
		internal.SaveAppScale()
	}

	now := time.Now()

	if now.Sub(r.lastCheckInternet) > r.checkInternetTimeout {
		if !content.Mgdata.ObjectPropExists("cache", "darkmode.data") {
			err := internal.SaveDarkMode()
			return err
		}

		go func() {
			if internal.IsInternetReachable() {
				content.IsInternet = true
			} else {
				content.IsInternet = false
			}
		}()
		r.lastCheckInternet = now
	}

	return nil
}

func (r *Root) Draw(context *guigui.Context, dst *ebiten.Image) {
	basicwidget.FillBackground(dst, context)
}
