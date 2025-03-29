/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Ilan Mayeux
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
	ESApp "eightysix/internal/app"
	"eightysix/internal/data"
	"eightysix/internal/file"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/rs/zerolog/log"
)

type Root struct {
	guigui.RootWidget

	lastCheckInternet    time.Time
	checkInternetTimeout time.Duration

	sidebar  Sidebar
	home     Home
	settings Settings
	//changelog Changelog
	//about     eightysix.About

	popup            basicwidget.Popup
	popupTitleText   basicwidget.Text
	popupCloseButton basicwidget.TextButton

	initOnce sync.Once
	err      error
}

func (r *Root) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	r.initOnce.Do(func() {
		app = &ESApp.App{
			FS:   &file.AppFS{GdataM: GDataM},
			Data: &data.Data{GDataM: GDataM},
		}

		r.checkInternetTimeout = time.Second
		if err := app.Data.InitColorMode(); err != nil {
			r.err = app.Error(err)
			return
		}
		if err := app.Data.InitAppScale(); err != nil {
			r.err = app.Error(err)
			return
		}
		log.Info().Msg("Init DarkMode and AppScale")
	})

	appender.AppendChildWidget(&r.sidebar)

	//u := float64(basicwidget.UnitSize(context))

	guigui.SetPosition(&r.sidebar, guigui.Position(r))
	sw, _ := r.sidebar.Size(context)
	p := guigui.Position(r)
	p.X += sw
	guigui.SetPosition(&r.home, p)
	guigui.SetPosition(&r.settings, p)

	switch r.sidebar.SelectedItemTag() {
	case "home":
		appender.AppendChildWidget(&r.home)
	case "settings":
		appender.AppendChildWidget(&r.settings)
	}

	// if len(content.Errs) != 0 {
	// 	r.popup.Open()
	// }
	// if len(content.Errs) > 0 {
	// 	contentWidth := int(12 * u)
	// 	contentHeight := int(6 * u)
	// 	bounds := guigui.Bounds(&r.popup)
	// 	contentPosition := image.Point{
	// 		X: bounds.Min.X + (bounds.Dx()-contentWidth)/2,
	// 		Y: bounds.Min.Y + (bounds.Dy()-contentHeight)/2,
	// 	}
	// 	contentBounds := image.Rectangle{
	// 		Min: contentPosition,
	// 		Max: contentPosition.Add(image.Pt(contentWidth, contentHeight)),
	// 	}
	// 	r.popup.SetContent(func(context *guigui.Context, appender *basicwidget.ContainerChildWidgetAppender) {
	// 		r.popupTitleText.SetText(content.Errs[0].Error())
	// 		r.popupTitleText.SetBold(true)
	// 		pt := contentBounds.Min.Add(image.Pt(int(0.5*u), int(0.5*u)))
	// 		guigui.SetPosition(&r.popupTitleText, pt)
	// 		appender.AppendChildWidget(&r.popupTitleText)
	//
	// 		r.popupCloseButton.SetText("Close")
	// 		r.popupCloseButton.SetOnUp(func() {
	// 			content.Errs = append(content.Errs[:0], content.Errs[1:]...)
	// 			r.popup.Close()
	// 		})
	// 		w, h := r.popupCloseButton.Size(context)
	// 		pt = contentBounds.Max.Add(image.Pt(-int(0.5*u)-w, -int(0.5*u)-h))
	// 		guigui.SetPosition(&r.popupCloseButton, pt)
	// 		appender.AppendChildWidget(&r.popupCloseButton)
	// 	})
	// 	r.popup.SetContentBounds(contentBounds)
	// 	r.popup.SetBackgroundBlurred(true)
	// 	r.popup.SetCloseByClickingOutside(false)
	//
	// 	appender.AppendChildWidget(&r.popup)
	// }
}

func (r *Root) Update(context *guigui.Context) error {
	if r.err != nil {
		return r.err
	}

	err := app.Data.UpdateData(context)
	if err != nil {
		return app.Error(err)
	}

	now := time.Now()

	if now.Sub(r.lastCheckInternet) > r.checkInternetTimeout {
		if !app.FS.IsDir() {
			err := app.Data.HandleDataReset()
			if err != nil {
				return app.Error(err)
			}
			log.Info().Msg("HandleDataReset")
		}

		go app.UpdateInternet()
		r.lastCheckInternet = now
	}

	return nil
}

func (r *Root) Draw(context *guigui.Context, dst *ebiten.Image) {
	basicwidget.FillBackground(dst, context)
}
