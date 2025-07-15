//go:build cjk
// +build cjk

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

package app

import (
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/basicwidget/cjkfont"
	"golang.org/x/text/language"
)

var localeItems = []basicwidget.DropdownListItem[language.Tag]{
	{
		Text: "English",
		ID:   language.English,
	},
}

func AppendRecommendedFaceSourceEntries(faceSourceEntries []basicwidget.FaceSourceEntry, locales []language.Tag) []basicwidget.FaceSourceEntry {
	return cjkfont.AppendRecommendedFaceSourceEntries(faceSourceEntries, locales)
}
