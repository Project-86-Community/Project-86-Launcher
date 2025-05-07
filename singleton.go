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
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"time"

	translator "github.com/Conight/go-googletrans"
	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type debugMode struct {
	Logs bool
}

type CacheT struct {
	Repo      *github.RepositoryRelease
	Timestamp time.Time
	ExpiresIn time.Duration
}

var (
	TheDebugMode debugMode
	ctx          = context.Background()
	githubClient = github.NewClient(nil)

	lBundle    *i18n.Bundle
	lLocalizer *i18n.Localizer
	t          = translator.New()

	e  = &debug.Debug{} // Debugger
	fs *file.AppFS      // Filesystem

	aErr *debug.Error
)

func GetAppError() *debug.Error {
	return aErr
}

func (r *Root) loadAndSetLocale(context *guigui.Context) *debug.Error {
	localePath := filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.LocaleFile)
	if err := fs.IsDir(e, localePath); err != nil {
		return r.model.data.SetLocale(context, language.English)
	}

	b, dErr := fs.Load(e, configs.Data, configs.LocaleFile)
	if dErr.Err != nil {
		return dErr
	}
	locale, err := language.Parse(string(b))
	if err != nil {
		return e.New(err, debug.FSError, debug.ErrDataLocaleLoad)
	}
	return r.model.data.SetLocale(context, locale)
}

func (r *Root) loadAndSetColorMode(context *guigui.Context) *debug.Error {
	colorModePath := filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.ColorModeFile)
	if err := fs.IsDir(e, colorModePath); err != nil {
		return r.model.data.SetColorMode(context, guigui.ColorModeLight)
	}

	b, dErr := fs.Load(e, configs.Data, configs.ColorModeFile)
	if dErr.Err != nil {
		return dErr
	}
	var colorMode guigui.ColorMode
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&colorMode); err != nil {
		return e.New(err, debug.FSError, debug.ErrDataColorModeLoad)
	}
	return r.model.data.SetColorMode(context, colorMode)
}

func (r *Root) loadAndSetAppScale(context *guigui.Context) *debug.Error {
	appScalePath := filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.AppScaleFile)
	if err := fs.IsDir(e, appScalePath); err != nil {
		return r.model.data.SetAppScale(context, 2)
	}

	b, dErr := fs.Load(e, configs.Data, configs.AppScaleFile)
	if dErr.Err != nil {
		return dErr
	}
	var appScale int
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&appScale); err != nil {
		return e.New(err, debug.FSError, debug.ErrDataAppScaleLoad)
	}
	return r.model.data.SetAppScale(context, appScale)
}

func (r *Root) loadAndSetCache() *debug.Error {
	cachePath := filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Cache, configs.CacheFile)
	if err := fs.IsDir(e, cachePath); err == nil {
		log.Info().Str("Cache", "cache found, loading cache...").Msg("Root.Init")
		b, dErr := fs.Load(e, configs.Cache, configs.CacheFile)
		if dErr.Err != nil {
			return dErr
		}
		var cache CacheT
		decoder := gob.NewDecoder(bytes.NewReader(b))
		if err := decoder.Decode(&cache); err != nil {
			return e.New(err, debug.FSError, debug.ErrCacheLoad)
		}
		return r.model.cache.SetCache(cache)
	}
	return nil
}
