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
	"image"
	"net/http"
	"p86l/assets"
	p86lLocale "p86l/assets/locale"
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Root struct {
	guigui.DefaultWidget

	border    basicwidget.Background
	bgImage   basicwidget.Image
	sidebar   Sidebar
	play      Play
	changelog Changelog
	settings  Settings
	about     About

	sync  sync.Once
	model Model
	dErr  *debug.Error
}

func (r *Root) IsCache() bool {
	cacheErr := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Cache, configs.CacheFile))
	if r.model.isInternet && (r.model.cache.repo == nil || cacheErr != nil) {
		return true
	}

	return false
}

func (r *Root) IsCacheOutdated() bool {
	if r.model.isInternet && r.model.cache.repo != nil && time.Since(r.model.cache.timestamp) > r.model.cache.expiresIn {
		return true
	}
	return false
}

func (r *Root) RunApp() *debug.Error {
	iconImages, dErr := assets.GetIconImages(e)
	if dErr != nil {
		return dErr
	}
	ebiten.SetWindowIcon(iconImages)

	afs, dErr := file.NewFS(e)
	if dErr.Err != nil {
		return dErr
	}
	fs = afs

	bundle, dErr := p86lLocale.GetLocales(e, language.English)
	if dErr != nil {
		return dErr
	}
	lBundle = bundle
	lLocalizer = i18n.NewLocalizer(bundle, "en")

	return nil
}

func (r *Root) Init(context *guigui.Context) *debug.Error {
	if err := r.loadAndSetLocale(context); err != nil {
		return err
	}

	if err := r.loadAndSetColorMode(context); err != nil {
		return err
	}

	if err := r.loadAndSetAppScale(context); err != nil {
		return err
	}

	if err := r.loadAndSetCache(); err != nil {
		return err
	}

	return nil
}

func (r *Root) CheckInternet() {
	check := func() bool {
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(configs.InternetServer)
		if err != nil {
			return false
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				return
			}
		}()

		return resp.StatusCode == 204
	}

	r.model.isInternet = check()
}

// func (r *Root) FetchRateLimitStatus() *debug.Error {
// 	rate, _, err := githubClient.RateLimit.Get(ctx)
// 	if err != nil {
// 		return e.New(err, debug.InternetError, debug.ErrRateLimit)
// 	}
// 	r.model.rateLimitTracker.Update(rate.Core.Remaining, rate.Core.Reset.Time)
// 	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
// }
//
// func (r *Root) FetchCache() *debug.Error {
// 	if r.model.rateLimitTracker.Valid() {
// 		repo, _, err := githubClient.Repositories.GetLatestRelease(ctx, configs.RepoOwner, configs.RepoName)
// 		if err != nil {
// 			return e.New(err, debug.InternetError, debug.ErrCacheInternet)
// 		}
// 		log.Info().Msg("Root.FetchCache")
//
// 		newCache := CacheT{
// 			Repo:      repo,
// 			Timestamp: time.Now(),
// 			ExpiresIn: time.Hour,
// 		}
// 		r.dErr = r.model.cache.SetCache(newCache)
// 	}
//
// 	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
// }

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		r.dErr = r.RunApp()
		if r.dErr != nil {
			return
		}
		r.dErr = r.Init(context)
		if r.dErr != nil {
			return
		}
	})

	if r.dErr != nil {
		aErr = r.dErr
		return aErr.Err
	}

	faceSources := []*text.GoTextFaceSource{
		basicwidget.DefaultFaceSource(),
	}
	for _, locale := range context.AppendLocales(nil) {
		fs := cjkfont.FaceSourceFromLocale(locale)
		if fs != nil {
			faceSources = append(faceSources, fs)
			break
		}
	}
	basicwidget.SetFaceSources(faceSources)

	borderBounds := context.Bounds(r)
	borderBounds.Min.X = context.Size(&r.sidebar).X - (basicwidget.UnitSize(context) / 12)
	borderBounds.Max.X = context.Size(&r.sidebar).X
	context.SetOpacity(&r.border, 0.7)

	img, dErr := assets.TheImageCache.Get(e, "banner")
	if dErr != nil {
		aErr = dErr
		return dErr.Err
	}
	r.bgImage.SetImage(img)
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	aspectRatio := float64(imgHeight) / float64(imgWidth)

	windowSize := context.Size(r)
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
	appender.AppendChildWidgetWithPosition(&r.bgImage, image.Pt(00, yOffset))

	r.play.SetModel(&r.model)
	r.changelog.SetModel(&r.model)
	r.settings.SetModel(&r.model)
	r.sidebar.SetModel(&r.model)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))
	{
		switch r.model.Mode() {
		case "play":
			appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
			appender.AppendChildWidgetWithBounds(&r.play, gl.CellBounds(1, 0))
		case "changelog":
			appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
			appender.AppendChildWidgetWithBounds(&r.changelog, gl.CellBounds(1, 0))
		case "settings":
			appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
			appender.AppendChildWidgetWithBounds(&r.settings, gl.CellBounds(1, 0))
		case "about":
			appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
			appender.AppendChildWidgetWithBounds(&r.about, gl.CellBounds(1, 0))
		}
	}

	return nil
}

func (r *Root) Update(context *guigui.Context) error {
	if ebiten.Tick()-r.model.root.lastInternetCheckTick >= int64(ebiten.TPS()) {
		if r.model.isInternet {
			// go func() {
			// 	dErr := r.FetchRateLimitStatus()
			// 	if dErr != nil {
			// 		e.SetToast(dErr)
			// 	}
			// }()
		}
		go r.CheckInternet()
		r.model.root.lastInternetCheckTick = ebiten.Tick()
	}

	if r.IsCache() {
		r.model.root.SetCacheDebounce(func() {
			log.Info().Str("Cache", "cache not found, downloading cache...").Msg("Root.Update")

			// dErr := r.FetchCache()
			// if dErr != nil {
			// 	e.SetToast(dErr)
			// }
		})
	}
	if r.IsCacheOutdated() {
		r.model.root.SetCacheDebounce(func() {
			log.Info().Str("Cache", "outdated, updating...").Msg("Root.Update")

			// dErr := r.FetchCache()
			// if dErr != nil {
			// 	e.SetToast(dErr)
			// }
		})
	}

	return nil
}
