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

package locale

import (
	"embed"
	"fmt"
	"io/fs"
	"p86l/internal/debug"
	"strings"

	"github.com/BurntSushi/toml"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locale.*.toml
var localeFS embed.FS

func GetLocales(appDebug *debug.Debug, locale language.Tag) (*i18n.Bundle, *debug.Error) {
	bundle := i18n.NewBundle(locale)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	entries, err := fs.ReadDir(localeFS, ".")
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFSFileNotExist)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".toml") {
			localeTag := extractLocaleTag(entry.Name())
			if localeTag != "" {
				_, err := bundle.LoadMessageFileFS(localeFS, entry.Name())
				if err != nil {
					return nil, appDebug.New(fmt.Errorf("%s: %v", entry.Name(), err), debug.FSError, debug.ErrFSFileNotExist)
				}
			}
		}
	}

	return bundle, nil
}

func extractLocaleTag(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) >= 2 && parts[0] == "locale" {
		return parts[1]
	}
	return ""
}
