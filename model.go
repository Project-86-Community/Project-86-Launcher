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
	pd "p86l/internal/debug"
	"p86l/internal/file"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/hashicorp/go-version"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Model struct {
	mode string

	progress string

	data  DataModel
	cache CacheModel
}

func (m *Model) Mode() string {
	if m.mode == "" {
		return "home"
	}
	return m.mode
}

func (m *Model) Progress() string {
	return m.progress
}

func (m *Model) SetMode(mode string) {
	log.Info().Str("Page", mode).Msg("Sidebar")
	m.mode = mode
}

func (m *Model) SetProgress(progress string) {
	m.progress = progress
}

func (m *Model) Data() *DataModel {
	return &m.data
}

func (m *Model) Cache() *CacheModel {
	return &m.cache
}

// -- DataModel: handles data for app --

type DataModel struct {
	validGameVersion bool
	file             file.Data
}

func NewData() file.Data {
	return file.Data{
		V:             0,
		Locale:        language.English.String(),
		AppScale:      2,
		ColorMode:     guigui.ColorModeLight,
		GameVersion:   "",
		UsePreRelease: false,
	}
}

func (d *DataModel) File() *file.Data {
	return &d.file
}

func (d *DataModel) IsValidGameVersion() bool {
	return d.validGameVersion
}

func (d *DataModel) GetAppScaleI(scale float64) int {
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

func (d *DataModel) SetLocale(context *guigui.Context, locale language.Tag) {
	log.Info().Any("Translation", locale).Str("DataModel", "SetLocale").Msg(pd.FileManager)
	d.file.Locale = locale.String()
	context.SetAppLocales([]language.Tag{locale})
	SetLanguage(locale.String())
}

func (d *DataModel) SetAppScale(context *guigui.Context, scale int) {
	log.Info().Any("Scaling", scale).Str("DataModel", "SetAppScale").Msg(pd.FileManager)
	d.file.AppScale = scale
	context.SetAppScale(d.GetAppScaleF(scale))
}

func (d *DataModel) SetColorMode(context *guigui.Context, mode guigui.ColorMode) {
	log.Info().Any("Theme", mode).Str("DataModel", "SetColorMode").Msg(pd.FileManager)
	d.file.ColorMode = mode
	context.SetColorMode(mode)
}

func (d *DataModel) SetUsePreRelease(value bool) *pd.Error {
	log.Info().Any("Pre-release", value).Str("DateModel", "SetUsePreRelease").Msg(pd.FileManager)
	d.file.UsePreRelease = value
	return nil
}

func (d *DataModel) SetGameVersion(ver string) *pd.Error {
	if ver == "" {
		return nil
	}

	_, err := version.NewVersion(ver)
	if err != nil {
		return E.New(err, pd.AppError, pd.ErrGameVersionInvalid)
	}

	log.Info().Any("Game Version", ver).Str("DateModel", "SetGameVersion").Msg(pd.FileManager)
	d.file.GameVersion = ver
	return nil
}

func (d *DataModel) SetFile(file file.Data) *pd.Error {
	log.Info().Str("DataModel", "SetFile").Msg(pd.FileManager)
	d.file = file
	return d.Save()
}

func (d *DataModel) Save() *pd.Error {
	log.Info().Str("DataModel", "Save").Msg(pd.FileManager)
	return SaveData(d.file)
}

// -- CacheModel --

type CacheModel struct {
	valid          bool
	file           file.Cache
	TranslatedBody string
}

func (c *CacheModel) File() *file.Cache {
	return &c.file
}

func (c *CacheModel) IsValid() bool {
	return c.valid
}

func (c *CacheModel) SetValid(value bool) {
	c.valid = value
}

func (c *CacheModel) SetRepo(repo *github.RepositoryRelease, locale string) *pd.Error {
	c.file.V = 0
	c.file.Repo = repo
	c.file.Timestamp = time.Now()
	c.file.ExpiresIn = time.Hour
	if err := c.file.Validate(E); err != nil {
		c.valid = false
	} else {
		c.valid = true
	}
	log.Info().Str("CacheModel", "SetRepo").Msg(pd.FileManager)
	return SaveCache(c.file)
}

func (c *CacheModel) Translate(locale string) {
	if !c.valid {
		return
	}

	if body := c.file.Repo.GetBody(); body != "" && locale != "en" {
		go func() {
			c.TranslatedBody = translateGT(body, locale)
		}()
	}
}
