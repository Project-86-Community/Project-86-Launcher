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

package main

import (
	"eightysix"
	"eightysix/content"
	"eightysix/content/icon"
	"eightysix/internal"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/quasilyte/gdata/v2"
)

type Root struct {
	guigui.RootWidget

	initOnce sync.Once

	lastCheckInternet    time.Time
	checkInternetTimeout time.Duration

	sidebar eightysix.Sidebar
	home    eightysix.Home
	//settings  eightysix.Settings
	//changelog eightysix.Changelog
	//about     eightysix.About
}

func (r *Root) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	r.initOnce.Do(func() {
		r.checkInternetTimeout = 2 * time.Second
	})

	appender.AppendChildWidget(&r.sidebar)

	guigui.SetPosition(&r.sidebar, guigui.Position(r))
	sw, _ := r.sidebar.Size(context)
	p := guigui.Position(r)
	p.X += sw
	guigui.SetPosition(&r.home, p)
	//guigui.SetPosition(&r.settings, p)
	//guigui.SetPosition(&r.changelog, p)
	//guigui.SetPosition(&r.about, p)

	switch r.sidebar.SelectedItemTag() {
	case "home":
		appender.AppendChildWidget(&r.home)
	case "settings":
		//appender.AppendChildWidget(&r.settings)
	case "changelog":
		//appender.AppendChildWidget(&r.changelog)
	case "about":
		//appender.AppendChildWidget(&r.about)
	}
}

func (r *Root) Update(context *guigui.Context) error {
	if content.Mgdata.ObjectPropExists("cache", "darkmode.data") {
		darkModeByte, err := content.Mgdata.LoadObjectProp("cache", "darkmode.data")
		if err != nil {
			return err
		}
		darkModeData, err := strconv.Atoi(string(darkModeByte))
		if err != nil {
			return err
		}
		context.SetColorMode(guigui.ColorMode(darkModeData))
	} else {
		darkModeData := guigui.ColorModeLight
		if err := content.Mgdata.SaveObjectProp("cache", "darkmode.data", []byte(fmt.Sprintf("%v", darkModeData))); err != nil {
			return err
		}
	}

	now := time.Now()

	if now.Sub(r.lastCheckInternet) > r.checkInternetTimeout {
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

func main() {
	var appName string
	if runtime.GOOS == "windows" {
		appName = "Project-86-Community\\Project-86-Launcher"
	} else {
		appName = "Project-86-Community/Project-86-Launcher"
	}

	m, err := gdata.Open(gdata.Config{
		AppName: appName,
	})
	if err != nil {
		panic(err)
	}
	content.Mgdata = m

	iconImages, err := icon.GetIconImages()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowIcon(iconImages)
	op := &guigui.RunOptions{
		Title:           "Project 86 Launcher",
		WindowMinWidth:  500,
		WindowMinHeight: 280,
	}
	if err = guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
