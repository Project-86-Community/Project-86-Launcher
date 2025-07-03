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

	"github.com/hajimehoshi/guigui"
	"github.com/hashicorp/go-version"
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

func (m *Model) Data() *DataModel {
	return &m.data
}

// -- DataModel: handles data for app --

type DataModel struct {
	validGameVersion bool
	file             file.Data
}

func NewData() file.Data {
	return file.Data{
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
	log.Info().Any("Translation", locale).Str("DataModel", "SetLocale").Msg("FileManager")
	d.file.Locale = locale.String()
	context.SetAppLocales([]language.Tag{locale})
	SetLanguage(locale.String())
}

func (d *DataModel) SetAppScale(context *guigui.Context, scale int) {
	log.Info().Any("Scaling", scale).Str("DataModel", "SetAppScale").Msg("FileManager")
	d.file.AppScale = scale
	context.SetAppScale(d.GetAppScaleF(scale))
}

func (d *DataModel) SetColorMode(context *guigui.Context, mode guigui.ColorMode) {
	log.Info().Any("Theme", mode).Str("DataModel", "SetColorMode").Msg("FileManager")
	d.file.ColorMode = mode
	context.SetColorMode(mode)
}

func (d *DataModel) SetUsePreRelease(value bool) *pd.Error {
	log.Info().Any("Pre-release", value).Str("DateModel", "SetUsePreRelease").Msg("FileManager")
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

	log.Info().Any("Game Version", ver).Str("DateModel", "SetGameVersion").Msg("FileManager")
	d.file.GameVersion = ver
	return nil
}

func (d *DataModel) Save() *pd.Error {
	log.Info().Str("DataModel", "Save").Msg("FileManager")
	return SaveData(d.file)
}

type CacheModel struct {
	file file.Cache
}
