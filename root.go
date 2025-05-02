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
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"net/http"
	"os"
	"p86l/assets"
	"p86l/assets/lang"
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Root struct {
	guigui.RootWidget

	border    basicwidget.Background
	bgImage   basicwidget.Image
	sidebar   Sidebar
	play      Play
	changelog Changelog
	settings  Settings
	about     About

	lastInternetCheckTick int64
	debounceCache         bool

	sync  sync.Once
	model Model
	dErr  *debug.Error
}

func (r *Root) RunApp() *debug.Error {
	e = &debug.Debug{}

	companyPath, dErr := file.GetCompanyPath(e)
	if dErr.Err != nil {
		return dErr
	}
	root, err := os.OpenRoot(companyPath)
	if err != nil {
		return e.New(err, debug.FSError, debug.ErrDirNotFound)
	}
	fs, dErr := file.NewFS(e, root, companyPath, configs.AppName)
	if dErr.Err != nil {
		return dErr
	}

	if TheDebugMode.IsRelease {
		logDir, dErr := fs.LogDir(e)
		if dErr.Err != nil {
			return dErr
		}

		if err := os.MkdirAll(logDir, 0755); err != nil {
			return e.New(err, debug.FSError, debug.ErrNewDirFailed)
		}

		timestamp := time.Now().Unix()
		logFileName := fmt.Sprintf("log_%d.log", timestamp)
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, err := os.Create(logFilePath)
		if err != nil {
			return e.New(err, debug.FSError, debug.ErrNewFileFailed)
		}

		TheDebugMode.LogFile = logFile

		multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	}

	if err := lang.GetLangs(); err != nil {
		return e.New(err, debug.FSError, debug.ErrFileNotFound)
	}

	r.CheckInternet()
	r.model.githubClient = github.NewClient(nil)

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (r *Root) Init(context *guigui.Context) *debug.Error {
	{
		value := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.LocaleFile))
		if value.Err != nil {
			dErr := r.model.data.SetLocale(context, language.English)
			if dErr.Err != nil {
				return dErr
			}
		}

		b, dErr := fs.Load(e, configs.Data, configs.LocaleFile)
		if dErr.Err != nil {
			return dErr
		}
		locale, err := language.Parse(string(b))
		if err != nil {
			return e.New(err, debug.FSError, debug.ErrLocaleLoad)
		}
		r.dErr = r.model.data.SetLocale(context, locale)
	}
	{
		value := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.ColorModeFile))
		if value.Err != nil {
			dErr := r.model.data.SetColorMode(context, guigui.ColorModeLight)
			if dErr.Err != nil {
				return dErr
			}
		}

		b, dErr := fs.Load(e, configs.Data, configs.ColorModeFile)
		if dErr.Err != nil {
			return dErr
		}
		var colorMode guigui.ColorMode
		decoder := gob.NewDecoder(bytes.NewReader(b))
		err := decoder.Decode(&colorMode)
		if err != nil {
			return e.New(err, debug.FSError, debug.ErrColorModeLoad)
		}
		r.dErr = r.model.data.SetColorMode(context, colorMode)
	}
	{
		value := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.AppScaleFile))
		if value.Err != nil {
			dErr := r.model.data.SetAppScale(context, 2)
			if dErr.Err != nil {
				return dErr
			}
		}

		b, dErr := fs.Load(e, configs.Data, configs.AppScaleFile)
		if dErr.Err != nil {
			return dErr
		}
		var appScale int
		decoder := gob.NewDecoder(bytes.NewReader(b))
		err := decoder.Decode(&appScale)
		if err != nil {
			return e.New(err, debug.FSError, debug.ErrAppScaleLoad)
		}
		r.dErr = r.model.data.SetAppScale(context, appScale)
	}

	{
		value := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Cache, configs.CacheFile))
		if value.Err == nil {
			log.Info().Str("Cache", "cache found, loading cache...").Msg("Root.Init")
			b, dErr := fs.Load(e, configs.Cache, configs.CacheFile)
			if dErr.Err != nil {
				return dErr
			}
			var cache CacheT
			decoder := gob.NewDecoder(bytes.NewReader(b))
			err := decoder.Decode(&cache)
			if err != nil {
				return e.New(err, debug.FSError, debug.ErrCacheLoad)
			}
			r.dErr = r.model.cache.SetCache(cache)
		}
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (r *Root) CheckInternet() {
	check := func() bool {
		client := http.Client{
			Timeout: 2 * time.Second,
		}

		resp, err := client.Get(configs.InternetServer)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == 204
	}

	r.model.isInternet = check()
}

func (r *Root) FetchRateLimitStatus() *debug.Error {
	rate, _, err := r.model.githubClient.RateLimit.Get(ctx)
	if err != nil {
		return e.New(err, debug.InternetError, debug.ErrRateLimit)
	}
	r.model.rateLimitTracker.Update(rate.Core.Remaining, rate.Core.Reset.Time)
	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (r *Root) FetchCache() *debug.Error {
	if r.model.rateLimitTracker.Valid() {
		repo, _, err := r.model.githubClient.Repositories.GetLatestRelease(ctx, configs.RepoOwner, configs.RepoName)
		if err != nil {
			return e.New(err, debug.InternetError, debug.ErrCacheInternet)
		}
		log.Info().Msg("Root.FetchCache")

		newCache := CacheT{
			Repo:      repo,
			Timestamp: time.Now(),
			ExpiresIn: time.Hour,
		}
		r.dErr = r.model.cache.SetCache(newCache)
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		r.dErr = r.RunApp()
		r.dErr = r.Init(context)
	})

	if r.dErr != nil && r.dErr.Err != nil {
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

	img, err := assets.TheImageCache.Get("banner")
	if err != nil {
		aErr = e.New(err, debug.FSError, debug.ErrFileNotFound)
		return err
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
	for i, bounds := range gl.CellBounds() {
		switch i {
		case 0:
			appender.AppendChildWidgetWithBounds(&r.sidebar, bounds)
		case 1:
			switch r.model.Mode() {
			case "play":
				appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
				appender.AppendChildWidgetWithBounds(&r.play, bounds)
			case "changelog":
				appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
				appender.AppendChildWidgetWithBounds(&r.changelog, bounds)
			case "settings":
				appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
				appender.AppendChildWidgetWithBounds(&r.settings, bounds)
			case "about":
				appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
				appender.AppendChildWidgetWithBounds(&r.about, bounds)
			}
		}
	}

	return nil
}

func (r *Root) Update(context *guigui.Context) error {
	if ebiten.Tick()-r.lastInternetCheckTick >= int64(ebiten.TPS()) {
		if r.model.isInternet {
			go func() {
				dErr := r.FetchRateLimitStatus()
				if dErr.Err != nil {
					e.SetToast(dErr)
				}
			}()
		}
		go r.CheckInternet()
		r.lastInternetCheckTick = ebiten.Tick()
	}

	if r.model.isInternet {
		value := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Cache, configs.CacheFile))
		if value.Err != nil {
			if r.model.isInternet && r.debounceCache {
				r.debounceCache = true
				go func() {
					log.Info().Str("Cache", "cache not found, downloading cache...").Msg("Root.Update")

					dErr := r.FetchCache()
					if dErr.Err != nil {
						e.SetToast(dErr)
					}
					r.debounceCache = false
				}()
			}
		}

		if r.model.cache.repo.GetBody() != "" {
			if time.Since(r.model.cache.timestamp) > r.model.cache.expiresIn && r.debounceCache {
				r.debounceCache = true
				go func() {
					log.Info().Str("Cache", "outdated, updating...").Msg("Root.Update")

					dErr := r.FetchCache()
					if dErr.Err != nil {
						e.SetToast(dErr)
					}
					r.debounceCache = false
				}()
			}
		}
	}

	return nil
}
