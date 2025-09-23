/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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

package images

import (
	"bytes"
	_ "embed"
	"image"
)

var (
	//go:embed icon.png
	icon []byte
	//go:embed banner.jpg
	banner []byte

	//go:embed buttons/ie.png
	iE []byte
	//go:embed buttons/github.png
	github []byte
	//go:embed buttons/discord.png
	discord []byte
	//go:embed buttons/patreon.png
	patreon []byte

	//go:embed about/lead.jpg
	leadDeveloper []byte
	//go:embed about/dev.jpg
	devDeveloper []byte
)

var (
	ImageIcon, ImageBanner, ImageIE, ImageGithub, ImageDiscord, ImagePatreon, ImageLeadDeveloper, ImageDevDeveloper image.Image
)

func init() {
	list := []struct {
		data   []byte
		target *image.Image
		name   string
	}{
		{data: icon, target: &ImageIcon, name: "icon"},
		{data: banner, target: &ImageBanner, name: "banner"},
		{data: iE, target: &ImageIE, name: "buttons/ie.png"},
		{data: github, target: &ImageGithub, name: "buttons/github.png"},
		{data: discord, target: &ImageDiscord, name: "buttons/discord.png"},
		{data: patreon, target: &ImagePatreon, name: "buttons/patreon.png"},
		{data: leadDeveloper, target: &ImageLeadDeveloper, name: "about/lead.jpg"},
		{data: devDeveloper, target: &ImageDevDeveloper, name: "about/dev.jpg"},
	}
	for _, img := range list {
		decoded, _, err := image.Decode(bytes.NewReader(img.data))
		if err != nil {
			panic(err)
		}
		*img.target = decoded
	}
}
