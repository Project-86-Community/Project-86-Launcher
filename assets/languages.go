/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2026 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
 * Copyright (C) 2026 Project 86 Community
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

package assets

import (
	"embed"
	"fmt"
	"maps"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed languages/*.toml
var languagesFS embed.FS

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	entries, err := languagesFS.ReadDir("languages")
	if err != nil {
		panic(fmt.Sprintf("i18n: failed to read languages dir: %v", err))
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := "languages/" + entry.Name()
		if _, err := bundle.LoadMessageFileFS(languagesFS, path); err != nil {
			panic(fmt.Sprintf("i18n: failed to load %s: %v", path, err))
		}
	}
}

// T is a self-contained localizer.
type T struct {
	loc *i18n.Localizer
}

// NewT creates a T for the given language tag (e.g. "en", "pt-BR").
// NewT falls back to English if the language is not found.
func NewT(lang string) T {
	return T{loc: i18n.NewLocalizer(bundle, lang, language.English.String())}
}

// Get translates a message by ID.
// Get returns the ID itself if not found,
func (t T) Get(id string, templateData ...map[string]any) string {
	cfg := &i18n.LocalizeConfig{MessageID: id}
	if len(templateData) > 0 {
		cfg.TemplateData = templateData[0]
	}
	msg, err := t.loc.Localize(cfg)
	if err != nil {
		return id
	}
	return msg
}

// Plural translates a pluralized message.
func (t T) Plural(id string, count int, templateData ...map[string]any) string {
	data := map[string]any{}
	if len(templateData) > 0 {
		maps.Copy(data, templateData[0])
	}
	data["Count"] = count

	msg, err := t.loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		PluralCount:  count,
		TemplateData: data,
	})
	if err != nil {
		return id
	}
	return msg
}
