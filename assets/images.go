// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

// Changed for p86l by realskyquest

package assets

import (
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"

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

func (i *imageCache) Get(name string) (*ebiten.Image, error) {
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
		defer f.Close()
		pImg, err = jpeg.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JPG: %w", err)
		}
	} else {
		// If JPG fails, try PNG
		f, err = pngImages.Open(name + ".png")
		if err != nil {
			return nil, fmt.Errorf("image not found as JPG or PNG: %w", err)
		}
		defer f.Close()
		pImg, err = png.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("failed to decode PNG: %w", err)
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
