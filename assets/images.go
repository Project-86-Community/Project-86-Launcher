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
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"p86l/internal/debug"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.jpg
var jpgImages embed.FS

//go:embed *.png
var pngImages embed.FS

type imageCacheKey struct {
	name string
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image
}

var TheImageCache = &imageCache{}

func (i *imageCache) Get(appDebug *debug.Debug, name string) (*ebiten.Image, *debug.Error) {
	key := imageCacheKey{name: name}

	// Check if image is already cached
	if img, ok := i.m[key]; ok {
		return img, nil
	}

	var (
		pImg image.Image
		err  error
		f    fs.File
	)

	// Try to open and decode as JPG
	f, err = jpgImages.Open(name + ".jpg")
	if err == nil {
		defer func() {
			if err := f.Close(); err != nil {
				return
			}
		}()
		pImg, err = jpeg.Decode(f)
		if err != nil {
			return nil, appDebug.New(fmt.Errorf("failed to decode JPG: %w", err), debug.FSError, debug.ErrFSFileNotExist)
		}
	} else {
		// If JPG fails, try PNG
		f, err = pngImages.Open(name + ".png")
		if err != nil {
			return nil, appDebug.New(fmt.Errorf("image not found as JPG or PNG: %w", err), debug.FSError, debug.ErrFSFileNotExist)
		}
		defer func() {
			if err := f.Close(); err != nil {
				return
			}
		}()
		pImg, err = png.Decode(f)
		if err != nil {
			return nil, appDebug.New(fmt.Errorf("failed to decode PNG: %w", err), debug.FSError, debug.ErrFSFileNotExist)
		}
	}

	img := ebiten.NewImageFromImage(pImg)

	// Initialize cache if nil
	if i.m == nil {
		i.m = make(map[imageCacheKey]*ebiten.Image)
	}
	i.m[key] = img
	return img, nil
}
