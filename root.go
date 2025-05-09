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
	"image"
	"p86l/assets"
	p86lLocale "p86l/assets/locale"
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/guigui/layout"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	border     basicwidget.Background
	bgImage    basicwidget.Image
	sidebar    Sidebar
	play       Play
	changelog  Changelog
	settings   Settings
	about      About

	sync  sync.Once
	model Model
	dErr  *debug.Error
}

func (r *Root) RunApp() *debug.Error {
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

func (r *Root) LoadB(context *guigui.Context, loadType, loadFile string) *debug.Error {
	// loadPath := filepath.Join(fs.DirAppPath(), loadDir, loadFile)
	// if dErr := fs.IsDir(e, loadPath); dErr != nil {
	// 	switch loadType {
	// 	case "locale":
	// 		return r.model.data.SetLocale(context, language.English)
	// 	case "appscale":
	// 		return r.model.data.SetColorMode(context, guigui.ColorModeLight)
	// 	case "colormode":
	// 		return r.model.data.SetAppScale(context, 2)
	// 	}
	// } else {
	// 	if loadDir == "cache" {
	// 		log.Info().Str("Cache", "cache found, loading cache...").Msg("Root.Init")
	// 		b, dErr := fs.Load(e, configs.Cache, configs.CacheFile)
	// 		if dErr != nil {
	// 			return dErr
	// 		}
	// 		var cache CacheT
	// 		decoder := gob.NewDecoder(bytes.NewReader(b))
	// 		if err := decoder.Decode(&cache); err != nil {
	// 			return e.New(err, debug.FSError, debug.ErrCacheLoad)
	// 		}
	// 		return r.model.cache.SetCache(cache)
	// 	}
	// }

	if dErr := fs.Stat(e, loadFile); dErr != nil {
		if loadType == "data" {
			var dataFile file.Data
			dataFile.Locale = language.English.String()
			dataFile.AppScale = 2
			dataFile.ColorMode = guigui.ColorModeLight
			return r.model.data.SetData(context, dataFile)
		}
	} else {

	}

	b, dErr := fs.Load(e, loadFile)
	if dErr != nil {
		return dErr
	}

	if loadType == "data" {
		var data file.Data
		decoder := gob.NewDecoder(bytes.NewReader(b))
		if err := decoder.Decode(&data); err != nil {
			return e.New(err, debug.FSError, debug.ErrDataLoad)
		}
		return r.model.data.SetData(context, data)
	}

	return nil
}

func (r *Root) Background(context *guigui.Context, appender *guigui.ChildWidgetAppender) *debug.Error {
	img, dErr := assets.TheImageCache.Get(e, "banner")
	if dErr != nil {
		return dErr
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

	return nil
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	r.sync.Do(func() {
		r.dErr = r.RunApp()
		r.dErr = r.LoadB(context, "data", configs.DataFile)
		//r.dErr = r.LoadB(context, "cache", configs.CacheFile)
	})

	if r.dErr != nil {
		gErr = r.dErr
		return r.dErr.Err
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

	dErr := r.Background(context, appender)
	if dErr != nil {
		gErr = dErr
		return dErr.Err
	}

	r.sidebar.SetModel(&r.model)
	r.play.SetModel(&r.model)
	r.changelog.SetModel(&r.model)
	r.settings.SetModel(&r.model)

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * basicwidget.UnitSize(context)),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))
	bounds := gl.CellBounds(1, 0)

	context.SetOpacity(&r.background, 0.6)
	context.SetOpacity(&r.border, 0.5)

	borderBounds := context.Bounds(r)
	borderBounds.Min.X = context.Size(&r.sidebar).X - (basicwidget.UnitSize(context) / 12)
	borderBounds.Max.X = context.Size(&r.sidebar).X

	if r.model.Mode() != "home" {
		appender.AppendChildWidgetWithBounds(&r.background, bounds)
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

	return nil
}
