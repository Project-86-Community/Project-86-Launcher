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
	"image"
	"p86l/assets"
	"p86l/internal/debug"
	"p86l/internal/widget"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Home struct {
	guigui.DefaultWidget

	vLayout         widget.VerticalLayout
	smLayoutForm    widget.Form
	smLayoutVLayout widget.VerticalLayout
	mdLayoutForm    widget.Form
	mdLayoutVLayout widget.VerticalLayout

	bannerImage basicwidget.Image
	titleText   basicwidget.Text
	gameButton  basicwidget.TextButton

	form          basicwidget.Form
	websiteButton basicwidget.TextButton
	githubButton  basicwidget.TextButton
	discordButton basicwidget.TextButton
	patreonButton basicwidget.TextButton

	err *debug.Error
}

func (h *Home) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	img, err := assets.TheImageCache.Get("banner")
	if err != nil {
		h.err = app.Debug.New(err, debug.FSError, debug.ErrFileNotFound)
		return
	}
	h.bannerImage.SetImage(img)

	h.websiteButton.SetOnDown(func() {
		// go func() {
		// 	if err := browser.OpenURL("https://taliayaya.github.io/Project-86-Website/"); err != nil {
		// 		app.PopError(errors.New(err.Error()))
		// 	}
		// }()
	})
	h.githubButton.SetOnDown(func() {
		// go func() {
		// 	if err := browser.OpenURL("https://github.com/Taliayaya/Project-86"); err != nil {
		// 		app.PopError(errors.New(err.Error()))
		// 	}
		// }()
	})
	h.discordButton.SetOnDown(func() {
		// go func() {
		// 	if err := browser.OpenURL("https://discord.gg/Yh2TQH97yA"); err != nil {
		// 		app.PopError(errors.New(err.Error()))
		// 	}
		// }()
	})
	h.patreonButton.SetOnDown(func() {
		// go func() {
		// 	if err := browser.OpenURL("https://patreon.com/project86"); err != nil {
		// 		app.PopError(errors.New(err.Error()))
		// 	}
		// }()
	})

	u := float64(basicwidget.UnitSize(context))
	w, _h := h.Size(context)
	pt := guigui.Position(h).Add(image.Pt(int(0.5*u), int(0.5*u)))

	{
		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()
		aspectRatio := float64(imgHeight) / float64(imgWidth)

		newWidth := w - int(u)
		newHeight := int(float64(newWidth) * aspectRatio)

		h.bannerImage.SetSize(context, newWidth+1, newHeight+1)
	}

	h.titleText.SetBold(true)
	h.titleText.SetText("Welcome to Project 86")
	h.titleText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

	h.websiteButton.SetText("Website")
	h.websiteButton.SetWidth(int(float64(w)/4) - int(1*u))

	h.githubButton.SetText("Github")
	h.githubButton.SetWidth(int(float64(w)/4) - int(1*u))

	h.discordButton.SetText("Discord")
	h.discordButton.SetWidth(int(float64(w)/4) - int(1*u))

	h.patreonButton.SetText("Patreon")
	h.patreonButton.SetWidth(int(float64(w)/4) - int(1*u))

	h.form.SetWidth(context, int(float64(w)/2)-int(0.5*u))
	h.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &h.websiteButton,
			SecondaryWidget: &h.githubButton,
		},
		{
			PrimaryWidget:   &h.discordButton,
			SecondaryWidget: &h.patreonButton,
		},
	})

	h.vLayout.SetHorizontalAlign(widget.HorizontalAlignCenter)
	h.vLayout.SetBackground(false)
	h.vLayout.SetLineBreak(false)
	h.vLayout.SetBorder(false)

	h.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&h.vLayout, pt)

	if app.IsInternet() {
		h.gameButton.SetText("Install")
		guigui.Enable(&h.gameButton)
	} else {
		h.gameButton.SetText("NO INTERNET")
		guigui.Disable(&h.gameButton)
	}

	if w >= int(940*context.AppScale()) {
		{
			imgWidth := img.Bounds().Dx()
			imgHeight := img.Bounds().Dy()
			aspectRatio := float64(imgHeight) / float64(imgWidth)

			newWidth := w - int(u)
			newHeight := int(float64(newWidth) * aspectRatio)

			h.bannerImage.SetSize(context, int(float64(newWidth)/1.8), int(float64(newHeight)/1.8))
		}
		h.gameButton.SetWidth(int(float64(w)/2.5) - int(1*u))
		h.titleText.ResetSize()
		h.titleText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

		h.titleText.SetScale(2.8)

		h.websiteButton.SetWidth(int(float64(w)/5) - int(1*u))
		h.githubButton.SetWidth(int(float64(w)/5) - int(1*u))
		h.discordButton.SetWidth(int(float64(w)/5) - int(1*u))
		h.patreonButton.SetWidth(int(float64(w)/5) - int(1*u))
		h.form.SetWidth(context, int(float64(w)/2.5)-int(0.5*u))

		h.mdLayoutVLayout.SetHorizontalAlign(widget.HorizontalAlignCenter)
		h.mdLayoutVLayout.SetBackground(false)
		h.mdLayoutVLayout.SetLineBreak(false)
		h.mdLayoutVLayout.SetBorder(false)

		h.mdLayoutVLayout.SetWidth(context, int(float64(w)/2.2)-int(2*u))
		h.mdLayoutVLayout.SetItems([]*widget.LayoutItem{
			{Widget: &h.titleText},
			{Widget: &h.gameButton},
			{Widget: &h.form},
		})

		h.mdLayoutForm.SetWidth(context, w-int(1*u))
		h.mdLayoutForm.SetItems([]*widget.FormItem{
			{PrimaryWidget: &h.bannerImage, SecondaryWidget: &h.mdLayoutVLayout},
		})

		_, mdLayoutFormHeight := h.mdLayoutForm.Size(context)
		guigui.SetPosition(&h.vLayout, image.Pt(pt.X, pt.Y+(_h/2-int(float64(mdLayoutFormHeight)/1.5))))
		h.vLayout.SetItems([]*widget.LayoutItem{
			{Widget: &h.mdLayoutForm},
		})
	} else if w >= int(640*context.AppScale()) {
		h.gameButton.SetWidth(int(float64(w)/2.3) - int(1*u))
		h.titleText.SetWidth(int(float64(w)/2.3) - int(1*u))
		h.titleText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

		h.titleText.SetScale(1.8)
		h.smLayoutVLayout.SetHorizontalAlign(widget.HorizontalAlignStart)
		h.smLayoutVLayout.SetBackground(false)
		h.smLayoutVLayout.SetLineBreak(false)
		h.smLayoutVLayout.SetBorder(false)

		h.smLayoutVLayout.SetWidth(context, w/2-int(2*u))
		h.smLayoutVLayout.SetItems([]*widget.LayoutItem{
			{Widget: &h.titleText},
			{Widget: &h.gameButton},
		})

		h.smLayoutForm.SetWidth(context, w-int(1*u))
		h.smLayoutForm.SetItems([]*widget.FormItem{
			{
				PrimaryWidget:   &h.smLayoutVLayout,
				SecondaryWidget: &h.form,
			},
		})

		h.vLayout.SetItems([]*widget.LayoutItem{
			{Widget: &h.bannerImage},
			{Widget: &h.smLayoutForm},
		})
	} else {
		h.gameButton.SetWidth(int(float64(w)/1.5) - int(1*u))
		h.titleText.ResetSize()

		h.titleText.SetScale(2)

		h.vLayout.SetItems([]*widget.LayoutItem{
			{Widget: &h.bannerImage},
			{Widget: &h.titleText},
			{Widget: &h.gameButton},
			{Widget: &h.form},
		})
	}
	appender.AppendChildWidget(&h.vLayout)
}

func (h *Home) Update(context *guigui.Context) error {
	if h.err != nil && h.err.Err != nil {
		AppErr = h.err
		return h.err.Err
	}
	return nil
}

func (h *Home) Size(context *guigui.Context) (int, int) {
	w, _h := guigui.Parent(h).Size(context)
	w -= sidebarWidth(context)
	return w, _h
}
