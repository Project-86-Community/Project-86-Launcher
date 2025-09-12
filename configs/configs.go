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

package configs

import "image"

var AppWindowMinSize = image.Pt(600, 300)

const (
	InternetServer = "https://clients3.google.com/generate_204"

	CompanyName = "Project-86-Community"
	AppName     = "Project-86-Launcher"
	AppTitle    = "Project 86 Launcher"

	RepoOwner = "Taliayaya"
	RepoName  = "Project-86"

	FolderLogs = "logs"
	FileData   = "data.json"
	FileCache  = "cache.json"

	Website = "https://project-86-community.github.io/Project-86-Website/"
	Github  = "https://github.com/Taliayaya/Project-86"
	Discord = "https://discord.gg/A8Fr6yEsUn"
	Patreon = "https://www.patreon.com/project86"
)
