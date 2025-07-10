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
	"os"
	pd "p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/guigui"
	"github.com/hashicorp/go-getter"
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

func IsValidGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		!strings.Contains(filename, "dev")
}

// -- downloading --

func DownloadGame(model *Model, src string) *pd.Error {
	model.SetProgress("Downloading...")
	buildDir := filepath.Join(FS.CompanyDirPath, "build")

	err := DownloadFile(model, src, buildDir, getter.ClientModeDir)
	if err != nil {
		return err
	}

	// Find the downloaded folder (assuming it's the only or newest folder in buildDir).
	entries, readErr := os.ReadDir(buildDir)
	if readErr != nil {
		return E.New(readErr, pd.FSError, pd.ErrFSDirRead)
	}

	// Find the folder that matches the pattern or just the most recently modified one.
	var downloadedFolder string
	var newestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			// Option 1: If you know it follows a pattern like "Project86-v*".
			if strings.HasPrefix(entry.Name(), "Project86-v") {
				downloadedFolder = entry.Name()
				break
			}

			// Option 2: Use the most recently modified folder.
			info, _ := entry.Info()
			if info.ModTime().After(newestTime) {
				newestTime = info.ModTime()
				downloadedFolder = entry.Name()
			}
		}
	}

	if downloadedFolder == "" {
		return E.New(fmt.Errorf("downloaded folder not found"), pd.FSError, pd.ErrFSDirNotExist)
	}

	// Rename the folder
	oldPath := filepath.Join(buildDir, downloadedFolder)
	newPath := filepath.Join(buildDir, "game")

	if renameErr := os.Rename(oldPath, newPath); renameErr != nil {
		return E.New(renameErr, pd.FSError, pd.ErrFSDirRename)
	}

	model.SetProgress("")

	return nil
}

func DownloadFile(model *Model, src, dest string, mode getter.ClientMode) *pd.Error {
	progressTracker := &DownloadProgressTracker{Model: model}
	client := &getter.Client{
		Src:              src,
		Dst:              dest,
		Mode:             mode,
		ProgressListener: progressTracker,
	}

	err := client.Get()
	if err != nil {
		return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
	}

	return nil
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
