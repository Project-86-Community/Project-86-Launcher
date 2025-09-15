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

package assets

import (
	_ "embed"
	"p86l/assets/images"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Icon          = ebiten.NewImageFromImage(images.ImageIcon)
	Banner        = ebiten.NewImageFromImage(images.ImageBanner)
	IE            = ebiten.NewImageFromImage(images.ImageIE)
	Github        = ebiten.NewImageFromImage(images.ImageGithub)
	Discord       = ebiten.NewImageFromImage(images.ImageDiscord)
	Patreon       = ebiten.NewImageFromImage(images.ImagePatreon)
	LeadDeveloper = ebiten.NewImageFromImage(images.ImageLeadDeveloper)
	DevDeveloper  = ebiten.NewImageFromImage(images.ImageDevDeveloper)

	//go:embed audio/p86l_ost_legion.ogg
	P86lOst []byte
)
