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
	"p86l/internal/file"

	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

type Model struct {
	mode string

	data  DataModel
	cache CacheModel
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
	dataFile file.Data
}

func (d *DataModel) File() file.Data {
	return d.dataFile
}

func (d *DataModel) SetLocale(context *guigui.Context, locale language.Tag) *debug.Error {
	log.Info().Any("Locale", locale).Msg("Model.DataModel.SetLocale")
	d.dataFile.Locale = locale.String()
	return d.SetData(context, d.dataFile)
}

func (d *DataModel) SetAppScale(context *guigui.Context, scale int) *debug.Error {
	log.Info().Any("Scale", scale).Msg("Model.DataModel.SetAppScale")
	d.dataFile.AppScale = scale
	return d.SetData(context, d.dataFile)
}

func (d *DataModel) SetColorMode(context *guigui.Context, mode guigui.ColorMode) *debug.Error {
	log.Info().Any("Mode", mode).Msg("Model.DataModel.SetColorMode")
	d.dataFile.ColorMode = mode
	return d.SetData(context, d.dataFile)
}

func (d *DataModel) SetData(context *guigui.Context, dataFile file.Data) *debug.Error {
	d.dataFile = dataFile

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&dataFile); err != nil {
		return e.New(err, debug.FSError, debug.ErrDataSave)
	}

	dErr := fs.Save(e, configs.DataFile, buf.Bytes())
	if dErr != nil {
		return dErr
	}

	locale, err := language.Parse(dataFile.Locale)
	if err != nil {
		return e.New(err, debug.DataError, debug.ErrDataLoad)
	}

	context.SetAppLocales([]language.Tag{locale})
	context.SetAppScale(d.GetAppScaleF(dataFile.AppScale))
	context.SetColorMode(dataFile.ColorMode)
	SetLanguage(dataFile.Locale)

	return nil
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
	cacheFile file.Cache

	isTranslate         bool
	translatedChangelog string
}

func (c *CacheModel) File() file.Cache {
	return c.cacheFile
}

func (c *CacheModel) IsTranslate() bool {
	return c.isTranslate
}

func (c *CacheModel) TranslatedChangelog() string {
	return c.translatedChangelog
}

func (c *CacheModel) SetCache(cacheFile file.Cache) *debug.Error {
	log.Info().Any("CacheModel.cache", cacheFile).Msg("SetCache")
	c.cacheFile = cacheFile

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(cacheFile); err != nil {
		return e.New(err, debug.FSError, debug.ErrCacheSave)
	}

	dErr := fs.Save(e, configs.CacheFile, buf.Bytes())
	if dErr != nil {
		return dErr
	}

	return nil
}

func (c *CacheModel) SetIsTranslate(value bool) {
	c.isTranslate = value
}

func (c *CacheModel) SetTranslatedChangelog(value string) {
	c.translatedChangelog = value
}
