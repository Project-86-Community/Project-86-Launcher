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
	"cmp"
	gctx "context"
	"image"
	"p86l"
	"p86l/assets"
	p86lLocale "p86l/assets/locale"
	"p86l/configs"
	pd "p86l/internal/debug"
	"p86l/internal/file"
	"runtime"
	"slices"
	"sync"
	"time"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
)

type Root struct {
	guigui.DefaultWidget

	background rootBackground
	sidebar    Sidebar
	play       Play
	changelog  Changelog
	settings   Settings
	about      About

	popup         basicwidget.Popup
	popupContent  rootPopupContent
	popupDebounce bool

	inProgress bool
	lastTick   int64
	model      p86l.Model

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry

	sync sync.Once
	err  *pd.Error
}

func (r *Root) runApp() *pd.Error {
	iconImages, err1 := assets.GetIconImages(p86l.E)
	afs, err2 := file.NewFS(p86l.E)
	bundle, err3 := p86lLocale.GetLocales(p86l.E, language.English)

	if err := cmp.Or(err1, err2, err3); err != nil {
		return err
	}

	ebiten.SetWindowIcon(iconImages)
	p86l.FS = afs
	p86l.LBundle = bundle
	p86l.LLocalizer = i18n.NewLocalizer(bundle, "en")

	return nil
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)

	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		err1 := r.runApp()
		err2 := p86l.LoadB(context, &r.model, "data")
		err3 := p86l.LoadB(context, &r.model, "cache")

		if err := cmp.Or(err1, err2, err3); err != nil {
			r.err = err
			return
		}

		r.model.Cache().Translate(r.model.Data().File().Locale)
		var gpuInfo ebiten.DebugInfo
		ebiten.ReadDebugInfo(&gpuInfo)

		log.Info().Msg("..:: GuiGui GUI Framework Alpha ::..")
		log.Info().Str("Version", p86l.TheDebugMode.Version).Msg("P86L - Project 86 Launcher")
		log.Info().Str("Detected OS", runtime.GOOS).Msg("Operating System")
		log.Info().Str("Graphics API", gpuInfo.GraphicsLibrary.String()).Msg("GPU")
		log.Warn().Str("LICENSE", p86l.ALicense).Msg("README")
	})

	if r.err != nil {
		p86l.GErr = r.err
		return r.err.Err
	}
	r.updateFontFaceSources(context)

	r.background.SetModel(&r.model)
	r.sidebar.SetModel(&r.model)
	r.play.SetModel(&r.model)
	r.changelog.SetModel(&r.model)
	r.settings.SetModel(&r.model)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * u),
			layout.FlexibleSize(1),
		},
	}
	r.background.SetSidebar(&r.sidebar)
	r.background.SetBgBounds(gl.CellBounds(1, 0))
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))

	switch r.model.Mode() {
	case "play":
		appender.AppendChildWidgetWithBounds(&r.play, gl.CellBounds(1, 0))
	case "changelog":
		appender.AppendChildWidgetWithBounds(&r.changelog, gl.CellBounds(1, 0))
	case "settings":
		appender.AppendChildWidgetWithBounds(&r.settings, gl.CellBounds(1, 0))
	case "about":
		appender.AppendChildWidgetWithBounds(&r.about, gl.CellBounds(1, 0))
	}

	// -- popup --

	r.popup.SetOnClosed(func(reason basicwidget.PopupClosedReason) {
		p86l.E.PopupErr = nil
		r.popupDebounce = false
	})
	if p86l.E.PopupErr != nil && !r.popupDebounce {
		r.popupDebounce = true
		r.popup.Open(context)
	}

	r.popupContent.popup = &r.popup
	r.popup.SetContent(&r.popupContent)
	r.popup.SetBackgroundBlurred(true)
	r.popup.SetCloseByClickingOutside(true)
	r.popup.SetAnimationDuringFade(false)

	appBounds := context.AppBounds()
	contentSize := image.Pt(int(12*u), int(6*u))
	popupPosition := image.Point{
		X: appBounds.Min.X + (appBounds.Dx()-contentSize.X)/2,
		Y: appBounds.Min.Y + (appBounds.Dy()-contentSize.Y)/2,
	}
	popupBounds := image.Rectangle{
		Min: popupPosition,
		Max: popupPosition.Add(contentSize),
	}
	context.SetSize(&r.popupContent, popupBounds.Size())
	appender.AppendChildWidgetWithBounds(&r.popup, popupBounds)

	return nil
}

func (r *Root) Tick(context *guigui.Context) error {
	if ebiten.Tick()-r.lastTick >= int64(ebiten.TPS()*5) && !r.inProgress {
		r.lastTick = ebiten.Tick()
		r.inProgress = true

		if cache := r.model.Cache(); !cache.IsValid() || time.Now().After(cache.File().Timestamp.Add(cache.File().ExpiresIn)) {
			log.Info().Str("Cache", "cache is invalid").Str("Root", "Tick").Msg(pd.NetworkManager)
			go func() {
				ctx := gctx.Background()
				release, _, rErr := p86l.GithubClient.Repositories.GetLatestRelease(ctx, configs.RepoOwner, configs.RepoName)
				if rErr != nil {
					log.Error().Any("Release", rErr).Msg(pd.NetworkManager)
					p86l.E.SetToast(p86l.E.New(rErr, pd.NetworkError, pd.ErrNetworkCacheRequest))
					r.inProgress = false
					return
				}
				err := cache.SetRepo(release, r.model.Data().File().Locale)
				if err != nil {
					p86l.E.SetToast(err)
				}
				r.inProgress = false
			}()
		}
	}

	return nil
}

type rootBackground struct {
	guigui.DefaultWidget

	bgImage    basicwidget.Image
	background basicwidget.Background

	model    *p86l.Model
	sidebar  *Sidebar
	bgBounds image.Rectangle

	err *pd.Error
}

func (r *rootBackground) SetModel(model *p86l.Model) {
	r.model = model
}

func (r *rootBackground) SetSidebar(sidebar *Sidebar) {
	r.sidebar = sidebar
}

func (r *rootBackground) SetBgBounds(bounds image.Rectangle) {
	r.bgBounds = bounds
}

func (r *rootBackground) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	img, err := assets.TheImageCache.Get(p86l.E, "banner")
	r.err = err

	if r.err != nil {
		p86l.GErr = r.err
		return r.err.Err
	}

	r.bgImage.SetImage(img)
	context.SetOpacity(&r.background, 0.9)

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	aspectRatio := float64(imgHeight) / float64(imgWidth)

	windowSize := context.ActualSize(r)
	availableWidth := windowSize.X

	newHeight := int(float64(availableWidth) * aspectRatio)

	if newHeight < windowSize.Y {
		newHeight = windowSize.Y
		availableWidth = int(float64(newHeight) / aspectRatio)
	}

	context.SetSize(&r.bgImage, image.Pt(availableWidth+2, newHeight+2))

	yOffset := 0
	if newHeight > windowSize.Y {
		yOffset = -(newHeight - windowSize.Y) / 2
	}

	imgPosition := image.Pt(00, yOffset)
	appender.AppendChildWidgetWithPosition(&r.bgImage, imgPosition)

	if r.model.Mode() != "home" {
		appender.AppendChildWidgetWithBounds(&r.background, r.bgBounds)
	}

	return nil
}

type rootPopupContent struct {
	guigui.DefaultWidget

	popup *basicwidget.Popup

	titleText   basicwidget.Text
	closeButton basicwidget.Button
}

func (r *rootPopupContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	r.titleText.SetValue(p86l.E.PopupErr.String())
	r.titleText.SetAutoWrap(true)
	r.titleText.SetBold(true)
	r.titleText.SetSelectable(true)

	r.closeButton.SetText("Close")
	r.closeButton.SetOnUp(func() {
		r.popup.Close()
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.LazySize(func(row int) layout.Size {
				if row != 1 {
					return layout.FixedSize(0)
				}
				return layout.FixedSize(r.closeButton.DefaultSize(context).Y)
			}),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.titleText, gl.CellBounds(0, 0))
	{
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(0, 1),
			Widths: []layout.Size{
				layout.FlexibleSize(1),
				layout.FixedSize(r.closeButton.DefaultSize(context).X),
			},
		}
		appender.AppendChildWidgetWithBounds(&r.closeButton, gl.CellBounds(1, 0))
	}

	return nil
}
