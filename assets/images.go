// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

// Changed for eightysix by realskyquest

package assets

import (
	"embed"
	"image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.jpg
var jngImages embed.FS

type imageCacheKey struct {
	name string
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image
}

var TheImageCache = &imageCache{}

func (i *imageCache) Get(name string) (*ebiten.Image, error) {
	key := imageCacheKey{
		name: name,
	}
	if img, ok := i.m[key]; ok {
		return img, nil
	}

	f, err := jngImages.Open(name + ".jpg")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	pImg, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	// if colorMode == guigui.ColorModeDark {
	// 	// Create a white image for dark mode.
	// 	rgbaImg := pImg.(draw.Image)
	// 	b := rgbaImg.Bounds()
	// 	for j := b.Min.Y; j < b.Max.Y; j++ {
	// 		for i := b.Min.X; i < b.Max.X; i++ {
	// 			if _, _, _, a := rgbaImg.At(i, j).RGBA(); a > 0 {
	// 				a16 := uint16(a)
	// 				rgbaImg.Set(i, j, color.RGBA64{a16, a16, a16, a16})
	// 			}
	// 		}
	// 	}
	// 	pImg = rgbaImg
	// }

	img := ebiten.NewImageFromImage(pImg)
	if i.m == nil {
		i.m = map[imageCacheKey]*ebiten.Image{}
	}
	i.m[key] = img
	return img, nil
}
