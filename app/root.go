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

package app

import (
	"image"
	"p86l"
	"p86l/assets"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"golang.org/x/text/language"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

type Root struct {
	guigui.DefaultWidget

	backgroundImage basicwidget.Image
	background      basicwidget.Background
	sidebar         Sidebar
	home            Home

	panelPlay     basicwidget.Panel
	panelSettings basicwidget.Panel
	panelAbout    basicwidget.Panel

	play     guigui.WidgetWithSize[*Play]
	settings guigui.WidgetWithSize[*Settings]
	about    guigui.WidgetWithSize[*About]

	model                   *p86l.Model
	backgroundImageSize     image.Point
	backgroundImagePosition image.Point

	sync sync.Once
}

func (r *Root) handleBackgroundImage(context *guigui.Context) {
	imgWidth := assets.Banner.Bounds().Dx()
	imgHeight := assets.Banner.Bounds().Dy()

	windowBounds := context.Bounds(r)
	windowWidth := windowBounds.Dx()
	windowHeight := windowBounds.Dy()

	imgAspectRatio := float64(imgWidth) / float64(imgHeight)
	windowAspectRatio := float64(windowWidth) / float64(windowHeight)

	var newWidth, newHeight int
	var xOffset, yOffset int

	if imgAspectRatio > windowAspectRatio {
		// The image is wider than the window. Scale by height and crop width.
		newHeight = windowHeight
		newWidth = int(float64(windowHeight) * imgAspectRatio)
		xOffset = 0
		yOffset = 0
	} else {
		// The image is taller than the window. Scale by width and crop height.
		newWidth = windowWidth
		newHeight = int(float64(windowWidth) / imgAspectRatio)
		xOffset = 0
		yOffset = (windowHeight - newHeight) / 2
	}

	r.backgroundImageSize = image.Pt(newWidth, newHeight)
	r.backgroundImagePosition = image.Pt(xOffset, yOffset)
}

func (r *Root) SetModel(model *p86l.Model) {
	r.model = model
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return r.model
	default:
		return nil
	}
}

func (r *Root) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&r.backgroundImage)
	adder.AddChild(&r.sidebar)
	if page := r.model.Data().Page(); page == p86l.PageHome {
		adder.AddChild(&r.home)
	} else {
		adder.AddChild(&r.background)
		switch page {
		case p86l.PagePlay:
			adder.AddChild(&r.panelPlay)
		case p86l.PageSettings:
			adder.AddChild(&r.panelSettings)
		case p86l.PageAbout:
			adder.AddChild(&r.panelAbout)
		}
	}
}

func (r *Root) Update(context *guigui.Context) error {
	r.sync.Do(func() {
		logger := r.model.Log().Logger()
		logger.Info().Str("Version", p86l.LauncherVersion).Msg("P86L - Project 86 Launcher")
		logger.Info().Str("Detected OS", runtime.GOOS).Msg("Operating System")

		var gpuInfo ebiten.DebugInfo
		ebiten.ReadDebugInfo(&gpuInfo)
		logger.Info().Str("Graphics API", gpuInfo.GraphicsLibrary.String()).Msg("GPU")

		data := r.model.Data()
		context.SetAppLocales([]language.Tag{data.Lang()})
		if data.IsNew() {
			colorMode := context.ColorMode()
			context.SetColorMode(colorMode)
		} else {
			if data.UseDarkmode() {
				context.SetColorMode(guigui.ColorModeDark)
			} else {
				context.SetColorMode(guigui.ColorModeLight)
			}
		}
		context.SetAppScale(data.AppScale())
	})

	r.backgroundImage.SetImage(assets.Banner)
	r.handleBackgroundImage(context)
	context.SetOpacity(&r.background, 0.9)

	u := basicwidget.UnitSize(context)
	x := context.Bounds(r).Size().X - (8 * u)

	{
		y := r.play.Widget().Overflow(context).Y
		r.play.SetFixedSize(image.Pt(x, y))
	}
	{
		y := r.settings.Widget().Overflow(context).Y
		r.settings.SetFixedSize(image.Pt(x, y))
	}
	{
		y := r.about.Widget().Overflow(context).Y
		r.about.SetFixedSize(image.Pt(x, y))
	}

	r.panelPlay.SetContent(&r.play)
	r.panelSettings.SetContent(&r.settings)
	r.panelAbout.SetContent(&r.about)

	return nil
}

func (r *Root) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &r.backgroundImage:
		return image.Rectangle{
			Min: r.backgroundImagePosition,
			Max: image.Pt(
				r.backgroundImagePosition.X+r.backgroundImageSize.X,
				r.backgroundImagePosition.Y+r.backgroundImageSize.Y,
			),
		}
	}

	u := basicwidget.UnitSize(context)
	layout := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &r.sidebar,
				Size:   guigui.FixedSize(8 * u),
			},
			{
				Size: guigui.FlexibleSize(1),
			},
		},
	}
	if widget == &r.sidebar {
		return layout.WidgetBounds(context, context.Bounds(r), widget)
	} else if widget == &r.background {
		return layout.ItemBounds(context, context.Bounds(r), 1)
	}
	return layout.ItemBounds(context, context.Bounds(r), 1)
}

func (r *Root) Close() error {
	return r.model.Close()
}
