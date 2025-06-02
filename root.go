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
	"context"
	"encoding/gob"
	"fmt"
	"image"
	"net/http"
	"p86l/assets"
	p86lLocale "p86l/assets/locale"
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/ebiten/v2"
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

	background basicwidget.Background
	border     basicwidget.Background
	bgImage    basicwidget.Image

	versionText basicwidget.Text
	infoBg      basicwidget.Background
	infoText    basicwidget.Text

	sidebar   Sidebar
	play      Play
	changelog Changelog
	settings  Settings
	about     About

	model Model

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry

	sync sync.Once
	dErr *debug.Error
}

// Reduce repeating `if (err) != nil` statements
func (r *Root) assertErr(dErr *debug.Error) {
	if dErr != nil {
		log.Error().Int("Code", dErr.Code).Any("Type", dErr.Type).Err((dErr.Err)).Str("Root", "assertErr").Msg("ErrorManager")
		r.dErr = dErr
	}
}

func (r *Root) checkInternet(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "HEAD", configs.InternetServer, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Warn().AnErr("Internet, resp.Body.Close()", err).Msg("Root.checkInternet")
		}
	}()

	return resp.StatusCode == 204
}

func (r *Root) checkGitHubRateLimit(ctx context.Context) *github.RateLimits {
	limits, _, err := githubClient.RateLimit.Get(ctx)
	if err != nil {
		return nil
	}
	return limits
}

func (r *Root) fetchCache() *github.RepositoryRelease {
	release, _, err := githubClient.Repositories.GetLatestRelease(ctx, configs.RepoOwner, configs.RepoName)
	if err != nil {
		return nil
	}
	return release
}

func (r *Root) runApp() *debug.Error {
	iconImages, dErr := assets.GetIconImages(e)
	if dErr != nil {
		return dErr
	}
	ebiten.SetWindowIcon(iconImages)

	afs, dErr := file.NewFS(e)
	if dErr != nil {
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

func (r *Root) loadB(context *guigui.Context, loadType, loadFile string) *debug.Error {
	if dErr := fs.Stat(e, loadFile); dErr != nil {
		switch loadType {
		case "data":
			log.Info().Str("Data", "data not found, creating data...").Str("Root", "loadB").Msg("FileManager")
			dataFile := NewData()
			dataFile.Log()
			return r.model.DataM.SetData(context, dataFile)
		case "cache":
			log.Info().Str("Cache", "cache not found").Str("Root", "loadB").Msg("FileManager")
			return nil
		}
	}

	b, dErr := fs.Load(e, loadFile)
	if dErr != nil {
		return dErr
	}
	decoder := gob.NewDecoder(bytes.NewReader(b))

	switch loadType {
	case "data":
		log.Info().Str("Data", "data found, loading data...").Str("Root", "loadB").Msg("FileManager")
		var dataFile file.Data

		if err := decoder.Decode(&dataFile); err != nil {
			return e.New(err, debug.FSError, debug.ErrDataLoad)
		}

		dataFile.Log()
		return r.model.DataM.SetData(context, dataFile)
	case "cache":
		log.Info().Str("Cache", "cache found, loading cache...").Str("Root", "loadB").Msg("FileManager")
		var cacheFile file.Cache

		if err := decoder.Decode(&cacheFile); err != nil {
			return e.New(err, debug.FSError, debug.ErrCacheLoad)
		}

		return r.model.CacheM.SetCache(cacheFile)
	}

	return nil
}

func (r *Root) buildBackground(context *guigui.Context, img *ebiten.Image) image.Point {
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

	return image.Pt(00, yOffset)
}

func (r *Root) buildVersionText(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	r.versionText.SetValue(TheDebugMode.Version)

	sidebarBounds := context.Bounds(&r.sidebar)
	textSize := context.Size(&r.versionText)

	xPos := sidebarBounds.Min.X + (sidebarBounds.Dx()-textSize.X)/2
	yPos := sidebarBounds.Max.Y - textSize.Y - 10

	context.SetOpacity(&r.versionText, 0.5)
	appender.AppendChildWidgetWithPosition(&r.versionText, image.Pt(xPos, yPos))
}

func (r *Root) buildInfo(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	status := "No connection"
	if r.model.networkState.InternetAvailable() {
		status = "Connected"
	}

	remaining := 0
	resetTime := ""
	if r.model.networkState.GitHubRateLimit() != nil {
		remaining = r.model.networkState.GitHubRateLimit().Core.Remaining
		resetTime = r.model.networkState.GitHubRateLimit().Core.Reset.Time.Format(time.DateTime)
	}

	r.infoText.SetScale(1)
	if remaining == 60 || !r.model.networkState.InternetAvailable() {
		r.infoText.SetValue(fmt.Sprintf("Internet: %s, API: %d/60", status, remaining))
	} else {
		r.infoText.SetValue(fmt.Sprintf("Internet: %s, API: %d/60 (reset until %s)", status, remaining, resetTime))
	}

	windowSize := context.Bounds(r).Size()
	textSize := context.Size(&r.infoText)

	bgWidth := textSize.X
	bgHeight := textSize.Y
	bgPos := image.Pt(
		windowSize.X-bgWidth,
		windowSize.Y-bgHeight,
	)
	textPos := image.Pt(
		bgPos.X,
		bgPos.Y,
	)

	if r.model.Mode() == "home" {
		appender.AppendChildWidgetWithBounds(&r.infoBg, image.Rect(
			bgPos.X,
			bgPos.Y,
			bgPos.X+bgWidth,
			bgPos.Y+bgHeight,
		))
	}
	appender.AppendChildWidgetWithPosition(&r.infoText, textPos)
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)

	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)

}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		r.assertErr(r.runApp())
		r.assertErr(r.loadB(context, "data", configs.DataFile))
		r.assertErr(r.loadB(context, "cache", configs.CacheFile))

		log.Info().Msg("..:: GuiGui GUI Framework Alpha ::..")
		log.Info().Str("Version", TheDebugMode.Version).Msg("P86L - Project 86 Launcher")
		log.Info().Str("Detected OS", runtime.GOOS).Msg("Operating System")
		log.Info().Str("Graphics API", "TODO:").Msg("GPU")
		log.Warn().Str("LICENSE", aLicense).Msg("README")
	})

	if r.dErr != nil {
		gErr = r.dErr
		return r.dErr.Err
	}

	img, dErr := assets.TheImageCache.Get(e, "banner")
	if dErr != nil {
		gErr = r.dErr
		return dErr.Err
	}
	r.bgImage.SetImage(img)
	imgPosition := r.buildBackground(context, img)
	appender.AppendChildWidgetWithPosition(&r.bgImage, imgPosition)

	r.sidebar.SetModel(&r.model)
	r.play.SetModel(&r.model)
	r.changelog.SetModel(&r.model)
	r.settings.SetModel(&r.model)

	r.updateFontFaceSources(context)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))

	// Background and sidebar border
	context.SetOpacity(&r.background, 0.9)
	context.SetOpacity(&r.border, 0.5)

	borderBounds := context.Bounds(r)
	borderBounds.Min.X = context.Size(&r.sidebar).X - (basicwidget.UnitSize(context) / 12)
	borderBounds.Max.X = context.Size(&r.sidebar).X

	if r.model.Mode() != "home" {
		appender.AppendChildWidgetWithBounds(&r.background, gl.CellBounds(1, 0))
		appender.AppendChildWidgetWithBounds(&r.border, borderBounds)
	}

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
	r.buildVersionText(context, appender)
	r.buildInfo(context, appender)

	return nil
}

// Enable Update status, when we have internet and valid cachefiles
func (r *Root) tickUpdate(context *guigui.Context) *debug.Error {
	dErr := fs.IsDir(e, fs.DirGamePath())
	if dErr != nil {
		r.model.lunch.Set(LunchStatusInstall)
		return nil
	}

	gameVersion := r.model.DataM.File().GameVersion
	if gameVersion != "" {
		isBig, err := IsNewVersion(gameVersion, r.model.CacheM.File().Repo.GetTagName())
		if err != nil {
			return e.New(err, debug.DataError, debug.ErrDataLoad)
		}
		if isBig {
			r.model.lunch.Set(LunchStatusUpdate)
		} else {
			r.model.lunch.Set(LunchStatusPlay)
		}
	} else {
		r.model.lunch.Set(LunchStatusPlay)
	}

	return nil
}

func (r *Root) Tick(_context *guigui.Context) error {
	r.model.networkState.netMutex.Lock()
	defer r.model.networkState.netMutex.Unlock()

	if ebiten.Tick()-r.model.networkState.lastCheckTick >= int64(ebiten.TPS()) && !r.model.networkState.checkingInProgress {
		r.model.networkState.checkingInProgress = true
		r.model.networkState.lastCheckTick = ebiten.Tick()

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			available := r.checkInternet(ctx)
			defer cancel()

			var limits *github.RateLimits
			if available {
				ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
				limits = r.checkGitHubRateLimit(ctx)
				cancel()

				dErr := fs.Stat(e, configs.CacheFile)
				if (dErr != nil && dErr.Code == debug.ErrFSRootFileNotExist) || !r.model.CacheM.IsVaild() {
					releases := r.fetchCache()
					if releases != nil {
						var cacheFile file.Cache
						cacheFile.Repo = releases
						cacheFile.Timestamp = time.Now()
						cacheFile.ExpiresIn = time.Hour
						r.model.CacheM.SetCache(cacheFile)
					}
				} else {
					dErr := r.tickUpdate(_context)
					if dErr != nil {
						e.SetToast(dErr)
					}
				}
			} else {
				// Enable play status, when build folder is detected and no internet
				dErr := fs.IsDir(e, fs.DirGamePath())
				if dErr == nil {
					r.model.lunch.Set(LunchStatusPlay)
				} else {
					r.model.lunch.Set(LunchStatusInstall)
				}
			}

			r.model.networkState.netMutex.Lock()
			defer r.model.networkState.netMutex.Unlock()
			r.model.networkState.internetAvailable = available
			r.model.networkState.githubRateLimit = limits
			r.model.networkState.checkingInProgress = false
		}()
	}

	return nil
}
