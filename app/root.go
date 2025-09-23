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

package app

import (
	"image"
	"os"
	"p86l"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/file"
	"p86l/internal/log"
	"runtime"
	"slices"
	"sync"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog"
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

	play     Play
	settings Settings
	about    About

	backgroundImageSize     image.Point
	backgroundImagePosition image.Point

	model *p86l.Model

	sync sync.Once

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry
}

func NewRoot(VERSION string) (*Root, *p86l.Model, *file.Filesystem, *zerolog.Logger, *os.File, error) {
	fs, err := file.NewFilesystem()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	logger, logFile, noFS, noAPI, err := log.NewLogger(VERSION, fs.Root())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	logger.Info().Str(log.Lifecycle, "app start").Msg(log.AppManager.String())
	logger.Info().Str(log.Lifecycle, "logging started").Msg(log.AppManager.String())
	logger.Info().Str("operating system", runtime.GOOS).Msg(log.AppManager.String())
	logger.Info().Str(log.Lifecycle, "init filesystem").Msg(log.FileManager.String())

	player, err := p86l.NewBGMPlayer()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	model := p86l.NewModel(logger, fs, player)

	if !noFS {
		dataSubModel := p86l.NewDataSubModel(model)
		model.AddSubModel(dataSubModel)
	}
	if !noAPI {
		cacheSubModel := p86l.NewCacheSubModel(model)
		model.AddSubModel(cacheSubModel)
	}

	model.Start()

	return &Root{model: model}, model, fs, logger, logFile, nil
}

func (r *Root) handleBackgroundImage(widgetBounds *guigui.WidgetBounds) {
	imgWidth := assets.Banner.Bounds().Dx()
	imgHeight := assets.Banner.Bounds().Dy()

	windowBounds := widgetBounds.Bounds()
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

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)

	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	//r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) contentWidget() guigui.Widget {
	dataFile := r.model.Data().Get()
	page := dataFile.Remember.Page

	switch p86l.SidebarPage(page) {
	case p86l.PageHome:
		return &r.home
	case p86l.PagePlay:
		return &r.play
	case p86l.PageSettings:
		return &r.settings
	case p86l.PageAbout:
		return &r.about
	}

	return nil
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return r.model
	default:
		return nil
	}
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	dataFile := r.model.Data().Get()
	page := dataFile.Remember.Page

	adder.AddChild(&r.backgroundImage)

	if p86l.SidebarPage(page) != p86l.PageHome {
		adder.AddChild(&r.background)
	}

	adder.AddChild(&r.sidebar)
	if content := r.contentWidget(); content != nil {
		adder.AddChild(content)
	}

	var err error
	r.sync.Do(func() {
		if context.ColorMode() == guigui.ColorModeDark {
			r.model.SetIsAutoUseDarkmode(true)
		}
		r.model.SetUIRefreshFn(func() {
			guigui.RequestRedraw(r)
		})
		r.model.SetSyncDataFn(func(m *p86l.Model, value bool) error {
			data := m.Data()
			cacheFile := m.Cache().Get()

			tag, err := data.Lang()
			if err != nil {
				return err
			}

			context.SetAppLocales([]language.Tag{tag})
			assets.LoadLanguage(tag.String())
			if !value {
				if data.UseDarkmode() {
					context.SetColorMode(guigui.ColorModeDark)
				} else {
					context.SetColorMode(guigui.ColorModeLight)
				}
			} else {
				if m.IsAutoUseDarkmode() {
					m.Data().Update(func(df *p86l.DataFile) {
						df.UseDarkmode = true
					})
				}
			}

			if cacheFile.Releases != nil && data.TranslateChangelog() && tag != language.English {
				m.Translate(p86l.ReleasesChangelogText(cacheFile, dataFile.UsePreRelease), tag.String())
			}

			context.SetAppScale(data.AppScale())

			remember := data.Remember()
			if remember.Active {
				if !value {
					ebiten.SetWindowSize(max(configs.AppWindowMinSize.X, remember.WSizeX), max(configs.AppWindowMinSize.Y, remember.WSizeY))
					ebiten.SetWindowPosition(max(0, remember.WPosX), max(0, remember.WPosY))
				}
				data.SetPage(p86l.SidebarPage(remember.Page))
			}

			if !data.DisableBgMusic() {
				m.BGMPlayer().Play()
			}
			return nil
		})
		err = r.model.SyncData()
	})
	if err != nil {
		return err
	}

	r.backgroundImage.SetImage(assets.Banner)
	{
		wMinX := int(float64(configs.AppWindowMinSize.X)*context.AppScale()) + basicwidget.UnitSize(context)*2
		wMinY := int(float64(configs.AppWindowMinSize.Y)*context.AppScale()) + basicwidget.UnitSize(context)*2
		ebiten.SetWindowSizeLimits(
			wMinX,
			wMinY,
			-1,
			-1,
		)
	}

	r.updateFontFaceSources(context)

	return nil
}

func (r *Root) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	data := r.model.Data()
	dataFile := data.Get()

	sx, sy := ebiten.WindowSize()
	px, py := ebiten.WindowPosition()
	page := dataFile.Remember.Page

	data.Update(func(df *p86l.DataFile) {
		df.Remember = p86l.DataRemember{
			WSizeX: sx,
			WSizeY: sy,
			WPosX:  px,
			WPosY:  py,
			Page:   int(page),
			Active: dataFile.Remember.Active,
		}
	})

	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)

	r.handleBackgroundImage(widgetBounds)
	bgBounds := image.Rect(
		r.backgroundImagePosition.X,
		r.backgroundImagePosition.Y,
		r.backgroundImagePosition.X+r.backgroundImageSize.X,
		r.backgroundImagePosition.Y+r.backgroundImageSize.Y,
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &r.backgroundImage,
				Size:   guigui.FlexibleSize(1),
			},
		},
	}).LayoutWidgets(context, bgBounds, layouter)

	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Size: guigui.FlexibleSize(1),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionHorizontal,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &r.sidebar,
							Size:   guigui.FixedSize(4 * u),
						},
						{
							Widget: &r.background,
							Size:   guigui.FlexibleSize(1),
							Layout: guigui.LinearLayout{
								Direction: guigui.LayoutDirectionVertical,
								Items: []guigui.LinearLayoutItem{
									{
										Widget: r.contentWidget(),
										Size:   guigui.FlexibleSize(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
