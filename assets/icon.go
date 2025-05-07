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
	"bytes"
	_ "embed"
	"image"
	"p86l/internal/debug"

	ico "github.com/biessek/golang-ico"
)

//go:embed p86l.ico
var p86lIcon []byte

func GetIconImages(appDebug *debug.Debug) ([]image.Image, *debug.Error) {
	var IconImages []image.Image

	reader := bytes.NewReader(p86lIcon)
	icons, err := ico.DecodeAll(reader)
	if err != nil {
		return nil, appDebug.New(err, debug.FSError, debug.ErrFSFileNotExist)
	}
	IconImages = append(IconImages, icons...)

	return IconImages, nil
}
