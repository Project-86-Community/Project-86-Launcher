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
	"fmt"
	"image"
	"p86l/internal/debug"
	"p86l/internal/widget"
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

	sidebar   Sidebar
	home      Home
	settings  Settings
	changelog Changelog
	about     About

	toast widget.Toast

	popup            basicwidget.Popup
	popupPanel       basicwidget.ScrollablePanel
	popupText        basicwidget.Text
	popupCloseButton basicwidget.TextButton

	initOnce sync.Once
	err      *debug.Error
}

func (r *Root) once() {
	r.checkInternetTimeout = time.Second
	// if err := Run(); err.Err != nil {
	// 	r.err = err
	// 	return
	// }
	if err := app.Data.InitColorMode(app.Debug); err.Err != nil {
		r.err = err
		return
	}
	if err := app.Data.InitAppScale(app.Debug); err.Err != nil {
		r.err = err
		return
	}
	log.Info().Msg("Init DarkMode and AppScale")
}

func (r *Root) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if ebiten.IsWindowBeingClosed() {
		log.Info().Msg("Closing App")
		if TheDebugMode.IsRelease {
			defer TheDebugMode.LogFile.Close()
		}
	}

	r.initOnce.Do(r.once)
	appender.AppendChildWidget(&r.sidebar)

	u := float64(basicwidget.UnitSize(context))
	w, h := r.Size(context)

	guigui.SetPosition(&r.sidebar, guigui.Position(r))
	sw, _ := r.sidebar.Size(context)
	p := guigui.Position(r)
	p.X += sw
	r.toast.SetWidth(context, w-sw)

	guigui.SetPosition(&r.home, p)
	guigui.SetPosition(&r.settings, p)
	guigui.SetPosition(&r.changelog, p)
	guigui.SetPosition(&r.about, p)
	guigui.SetPosition(&r.toast, p.Add(image.Pt(0, h-int(1.5*u))))

	switch r.sidebar.SelectedItemTag() {
	case "home":
		appender.AppendChildWidget(&r.home)
	case "settings":
		appender.AppendChildWidget(&r.settings)
	case "changelog":
		appender.AppendChildWidget(&r.changelog)
	case "about":
		appender.AppendChildWidget(&r.about)
	}

	if app.Debug.ToastErr != nil && app.Debug.ToastErr.Err != nil {
		r.toast.SetText(fmt.Sprintf("Code: %d, Type: %s, Error: %s", app.Debug.ToastErr.Code, string(app.Debug.ToastErr.Type), app.Debug.ToastErr.Err.Error()))
		r.toast.SetOnDown(func() {
			app.Debug.ToastErr = app.Debug.New(nil, debug.UnknownError, debug.ErrUnknown)
		})
		appender.AppendChildWidget(&r.toast)
	}

	// if len(app.Errs) != 0 {
	// 	r.popup.Open()
	// }
	// if len(app.Errs) > 0 {
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
	// 		r.popupPanel.SetSize(context, contentWidth-int(u), contentHeight-int(2*u))
	//
	// 		r.popupText.SetBold(true)
	//
	// 		log.Debug().Errs("TEST", app.Errs).Send()
	//
	// 		var popupMessage string
	// 		for _, appErr := range app.Errs {
	// 			popupMessage += "\n" + RemoveLineBreaks(appErr.Error())
	// 		}
	//
	// 		r.popupText.SetText(popupMessage)
	//
	// 		r.popupPanel.SetContent(func(context *guigui.Context, childAppender *basicwidget.ContainerChildWidgetAppender, offsetX, offsetY float64) {
	// 			p := guigui.Position(&r.popupPanel).Add(image.Pt(int(offsetX), int(offsetY)))
	//
	// 			guigui.SetPosition(&r.popupText, image.Pt(p.X, p.Y))
	// 			childAppender.AppendChildWidget(&r.popupText)
	// 		})
	//
	// 		pt := contentBounds.Min.Add(image.Pt(int(0.5*u), int(0.5*u)))
	//
	// 		guigui.SetPosition(&r.popupPanel, pt)
	// 		appender.AppendChildWidget(&r.popupPanel)
	//
	// 		r.popupCloseButton.SetText("Close All")
	// 		r.popupCloseButton.SetOnUp(func() {
	// 			app.Errs = app.Errs[:0]
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
	if r.err != nil && r.err.Err != nil {
		AppErr = r.err
		return r.err.Err
	}

	err := app.Data.UpdateData(context, app.Debug)
	if err.Err != nil {
		AppErr = err
		return err.Err
	}

	now := time.Now()

	if now.Sub(r.lastCheckInternet) > r.checkInternetTimeout {
		if !app.FS.IsDir() {
			err := app.Data.HandleDataReset(app.Debug)
			if err.Err != nil {
				AppErr = err
				return err.Err
			}
			log.Info().Msg("HandleDataReset")
		}

		go app.UpdateInternet()
		r.lastCheckInternet = now

		//app.PopupError(errors.New("YES"))
		//app.PopupError(errors.New("YEAKPWOKDPWKDPOWKDPOWKDPOWKDPOKWDOWKDPOWKDOWKDPWKDOPWKDPOWKDOPWKDPOWKDPOWKDPOWKDPOWKDPOWDKPWOKDOPWDKPWODKWPODKPOWDKS"))
	}

	//app.Update(githubClient, githubContext)

	return nil
}

func (r *Root) Draw(context *guigui.Context, dst *ebiten.Image) {
	basicwidget.FillBackground(dst, context)
}
