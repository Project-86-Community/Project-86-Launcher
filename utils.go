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
	"fmt"
	pd "p86l/internal/debug"
	"p86l/internal/file"

	"github.com/hajimehoshi/guigui"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

func LoadB(context *guigui.Context, model *Model, loadType string) *pd.Error {
	switch loadType {
	case "data":
		if err := FS.IsDirR(E, FS.FileDataPath()); err != nil {
			log.Info().Str("Data", "data not found, creating data...").Str("utils", "loadB").Msg(pd.FileManager)
			d := NewData()
			d.Log()
			model.data.file = d
			return model.data.Save()
		}
	case "cache":
		if err := FS.IsDirR(E, FS.FileCachePath()); err != nil {
			log.Info().Str("Cache", "cache not found").Str("utils", "loadB").Msg(pd.FileManager)
			return nil
		}
	}

	switch loadType {
	case "data":
		d, err := LoadData()
		if err != nil {
			return err
		}

		tag, rErr := language.Parse(d.Locale)
		if rErr != nil {
			return E.New(rErr, pd.DataError, pd.ErrDataLocaleInvalid)
		}
		model.data.SetLocale(context, tag)
		model.data.SetAppScale(context, d.AppScale)
		model.data.SetColorMode(context, d.ColorMode)
		model.data.SetUsePreRelease(d.UsePreRelease)
		return model.data.SetGameVersion(d.GameVersion)
	case "cache":
		c, err := LoadCache()
		if err != nil {
			return err
		}
		if err := c.Validate(E); err == nil {
			model.cache.valid = true
		}
		model.cache.file = *c
	}
	return nil
}

func SetLanguage(lang string) {
	LLocalizer = i18n.NewLocalizer(LBundle, lang)
}

func T(key string) string {
	lMsg, err := LLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return fmt.Sprintf("!{%s}", key)
	}

	return lMsg
}

func translateGT(body string, target string) string {
	result, err := t.Translate(body, "auto", target)
	if err != nil {
		return "?"
	}
	return result.Text
}

func OpenBrowser(url string) {
	log.Info().Str("Url", url).Msg("OpenBrowser")
	if err := browser.OpenURL(url); err != nil {
		E.SetPopup(E.New(err, pd.AppError, pd.ErrBrowserOpen))
	}
}

// -- Funcs for loading and saving --

func LoadData() (*file.Data, *pd.Error) {
	b, err := FS.Load(E, FS.FileDataPath())
	if err != nil {
		return nil, err
	}

	d, err := FS.DecodeData(E, b)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func LoadCache() (*file.Cache, *pd.Error) {
	b, err := FS.Load(E, FS.FileCachePath())
	if err != nil {
		return nil, err
	}

	c, err := FS.DecodeCache(E, b)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func SaveData(d file.Data) *pd.Error {
	b, err := FS.EncodeData(E, d)
	if err != nil {
		return err
	}

	err = FS.Save(E, FS.FileDataPath(), b)
	if err != nil {
		return err
	}

	return nil
}

func SaveCache(c file.Cache) *pd.Error {
	b, err := FS.EncodeCache(E, c)
	if err != nil {
		return err
	}

	err = FS.Save(E, FS.FileCachePath(), b)
	if err != nil {
		return err
	}

	return nil
}
