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
	"p86l/configs"
	"p86l/internal/debug"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Model struct {
	mode             string
	isInternet       bool
	rateLimitTracker RateLimitTracker
	data             DataModel
	cache            CacheModel
}

func (m *Model) Mode() string {
	if m.mode == "" {
		return "home"
	}
	return m.mode
}

func (m *Model) SetMode(mode string) {
	log.Info().Str("Page", mode).Msg("Sidebar")
	m.mode = mode
}

type DataModel struct {
	locale    language.Tag
	colorMode guigui.ColorMode
	appScale  int
}

func (d *DataModel) SetLocale(context *guigui.Context, locale language.Tag) *debug.Error {
	log.Info().Str("Locale", locale.String()).Msg("SetLocale")
	d.locale = locale

	context.SetAppLocales([]language.Tag{d.locale})
	SetLanguage(locale.String())

	dErr := fs.Save(e, configs.Data, configs.LocaleFile, []byte(locale.String()))
	if dErr.Err != nil {
		return dErr
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *DataModel) SetColorMode(context *guigui.Context, colorMode guigui.ColorMode) *debug.Error {
	log.Info().Int("guigui.ColorMode", int(colorMode)).Msg("SetColorMode")
	d.colorMode = colorMode
	context.SetColorMode(d.colorMode)

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(colorMode); err != nil {
		return e.New(err, debug.FSError, debug.ErrColorModeSave)
	}

	dErr := fs.Save(e, configs.Data, configs.ColorModeFile, buf.Bytes())
	if dErr.Err != nil {
		return dErr
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *DataModel) SetAppScale(context *guigui.Context, scale int) *debug.Error {
	log.Info().Float64("guigui.AppScale", d.GetAppScaleF(scale)).Msg("SetAppScale")
	d.appScale = scale
	context.SetAppScale(d.GetAppScaleF(d.appScale))

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(scale); err != nil {
		return e.New(err, debug.FSError, debug.ErrAppScaleSave)
	}

	dErr := fs.Save(e, configs.Data, configs.AppScaleFile, buf.Bytes())
	if dErr.Err != nil {
		return dErr
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *DataModel) GetAppScale(scale float64) int {
	switch scale {
	case 0.5: // 50%
		return 0
	case 0.75: // 75%
		return 1
	case 1.0: // 100%
		return 2
	case 1.25: // 125%
		return 3
	case 1.50: // 150%
		return 4
	}

	return -1
}

func (d *DataModel) GetAppScaleF(scale int) float64 {
	switch scale {
	case 0: // 50%
		return 0.5
	case 1: // 75%
		return 0.75
	case 2: // 100%
		return 1.0
	case 3: // 125%
		return 1.25
	case 4: // 150%
		return 1.50
	}

	return -1
}

type CacheModel struct {
	repo                *github.RepositoryRelease
	timestamp           time.Time
	expiresIn           time.Duration
	translatedChangelog string
}

func (c *CacheModel) SetCache(cache CacheT) *debug.Error {
	log.Info().Any("CacheModel.cache", cache).Msg("SetCache")
	c.repo = cache.Repo
	c.timestamp = cache.Timestamp
	c.expiresIn = cache.ExpiresIn

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(cache); err != nil {
		return e.New(err, debug.FSError, debug.ErrCacheSave)
	}

	dErr := fs.Save(e, configs.Cache, configs.CacheFile, buf.Bytes())
	if dErr.Err != nil {
		return dErr
	}

	return e.New(nil, debug.UnknownError, debug.ErrUnknown)
}
