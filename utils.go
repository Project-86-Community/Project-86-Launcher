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
	"bytes"
	"fmt"
	"image"
	"os"
	"p86l/assets"
	"runtime"
	"strings"

	"github.com/fyne-io/image/ico"
)

var DisableAPI bool = false

func GetIcons() ([]image.Image, error) {
	images, err := ico.DecodeAll(bytes.NewReader(assets.P86lIco))
	if err != nil {
		return nil, fmt.Errorf("failed to decode icons: %w", err)
	}

	return images, nil
}

func GetUsername() string {
	var username string
	switch runtime.GOOS {
	case "windows":
		username = os.Getenv("USERNAME")
	default:
		username = os.Getenv("USER")
	}

	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	return strings.TrimSpace(username)
}

func MergeRectangles(rects ...image.Rectangle) image.Rectangle {
	if len(rects) == 0 {
		return image.Rectangle{}
	}

	mergedMin := rects[0].Min
	mergedMax := rects[0].Max

	for _, r := range rects[1:] {
		if r.Min.X < mergedMin.X {
			mergedMin.X = r.Min.X
		}
		if r.Min.Y < mergedMin.Y {
			mergedMin.Y = r.Min.Y
		}
		if r.Max.X > mergedMax.X {
			mergedMax.X = r.Max.X
		}
		if r.Max.Y > mergedMax.Y {
			mergedMax.Y = r.Max.Y
		}
	}

	return image.Rectangle{
		Min: mergedMin,
		Max: mergedMax,
	}
}
