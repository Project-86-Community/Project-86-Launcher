// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package content

import (
	"embed"
	// eightysix Change Start
	"image/jpeg"
	// eightysix Change End

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
)

// eightysix Change Start

//go:embed *.jpg
var jngImages embed.FS

// eightysix Change End

type imageCacheKey struct {
	name      string
	colorMode guigui.ColorMode
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image
}

var TheImageCache = &imageCache{}

func (i *imageCache) Get(name string, colorMode guigui.ColorMode) (*ebiten.Image, error) {
	key := imageCacheKey{
		name:      name,
		colorMode: colorMode,
	}
	if img, ok := i.m[key]; ok {
		return img, nil
	}

	// eightysix Change Start
	f, err := jngImages.Open(name + ".jpg")
	// eightysix Change End
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// eightysix Change Start
	pImg, err := jpeg.Decode(f)
	// eightysix Change End
	if err != nil {
		return nil, err
	}

	// eightysix Change Start

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

	// eightysix Change End

	img := ebiten.NewImageFromImage(pImg)
	if i.m == nil {
		i.m = map[imageCacheKey]*ebiten.Image{}
	}
	i.m[key] = img
	return img, nil
}
