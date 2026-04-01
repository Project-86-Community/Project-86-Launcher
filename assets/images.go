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
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/fyne-io/image/ico"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images/*
var imagesFS embed.FS

var (
	// Images holds all jpg/png files as *ebiten.Image, keyed by filename without extension.
	// e.g. "images/logo.png" → Images["logo"]
	Images map[string]*ebiten.Image

	// Icons holds all ico files as []image.Image.
	// e.g. "images/app.ico" → Icons["app"]
	Icons map[string][]image.Image
)

func init() {
	Images = make(map[string]*ebiten.Image)
	Icons = make(map[string][]image.Image)

	entries, err := imagesFS.ReadDir("images")
	if err != nil {
		panic(fmt.Sprintf("assets: failed to read images dir: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		key := strings.TrimSuffix(name, filepath.Ext(name))
		path := "images/" + name

		data, err := imagesFS.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("assets: failed to read %s: %v", path, err))
		}

		switch ext {
		case ".jpg", ".jpeg", ".png":
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(fmt.Sprintf("assets: failed to decode image %s: %v", path, err))
			}
			Images[key] = ebiten.NewImageFromImage(img)

		case ".ico":
			imgs, err := ico.DecodeAll(bytes.NewReader(data))
			if err != nil {
				panic(fmt.Sprintf("assets: failed to decode ico %s: %v", path, err))
			}
			Icons[key] = imgs
		}
	}
}
