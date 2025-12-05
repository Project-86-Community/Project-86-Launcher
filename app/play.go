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
	"p86l"
	"p86l/assets"
	"p86l/configs"
	"time"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Play struct {
	guigui.DefaultWidget

	actionButtons                                                   [3]basicwidget.Button
	form                                                            basicwidget.Form
	gameVersionText, versionText, downloadsText, totalDownloadsText basicwidget.Text
	prereleaseText                                                  basicwidget.Text
	prereleaseToggle                                                basicwidget.Toggle
	changelogPanel                                                  basicwidget.Panel
	changelogText                                                   basicwidget.Text
	linkButtons                                                     [4]basicwidget.Button
}

func (p *Play) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	for i := range p.actionButtons {
		adder.AddChild(&p.actionButtons[i])
	}
	adder.AddChild(&p.form)
	adder.AddChild(&p.changelogPanel)
	for i := range p.linkButtons {
		adder.AddChild(&p.linkButtons[i])
	}

	model := context.Model(p, modelKeyModel).(*p86l.Model)
	data := model.Data()
	dataFile := data.Get()
	cacheFile := model.Cache().Get()

	inProgress := model.InProgress()
	if inProgress {
		for i := range p.actionButtons {
			context.SetEnabled(&p.actionButtons[i], !inProgress)
		}
	} else {
		// Install & Play
		if dataFile.UsePreRelease {
			preReleaseAvail, _ := model.CheckFilesCached(p86l.PathGamePreRelease)
			context.SetEnabled(&p.actionButtons[0], !preReleaseAvail)
			context.SetEnabled(&p.actionButtons[2], preReleaseAvail)
		} else {
			stableAvail, _ := model.CheckFilesCached(p86l.PathGameStable)
			context.SetEnabled(&p.actionButtons[0], !stableAvail)
			context.SetEnabled(&p.actionButtons[2], stableAvail)
		}
	}

	actionTexts := [3]string{p86l.T("play.install"), p86l.T("play.update"), p86l.T("play.play")}
	for i := range p.actionButtons {
		p.actionButtons[i].SetText(actionTexts[i])
	}

	p.actionButtons[0].SetOnDown(func() {
		go func() {
			model.InProgress(true)
			time.Sleep(time.Second * 5)
			model.InProgress(false)
			for i := range p.actionButtons {
				guigui.RequestRedraw(&p.actionButtons[i])
			}
		}()
	})

	p.changelogText.SetAutoWrap(true)
	p.changelogText.SetMultiline(true)
	if cacheFile.Releases != nil && dataFile.TranslateChangelog && dataFile.Lang != "en" {
		p.changelogText.SetValue(cacheFile.ChangelogTranslation)
	} else {
		p.changelogText.SetValue(p86l.ReleasesChangelogText(cacheFile, dataFile.UsePreRelease))
	}

	p.changelogPanel.SetContent(&p.changelogText)
	p.changelogPanel.SetAutoBorder(true)
	p.changelogPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	p.gameVersionText.SetValue(p86l.T("play.version"))
	p.downloadsText.SetValue(p86l.T("play.total"))
	p.prereleaseText.SetValue(p86l.T("play.prerelease"))

	p.versionText.SetValue(p86l.GameVersionText(cacheFile, dataFile.UsePreRelease))
	p.totalDownloadsText.SetValue(p86l.ReleasesDownloadCountText(cacheFile, dataFile.UsePreRelease))

	p.prereleaseToggle.SetOnValueChanged(func(value bool) {
		data.Update(func(df *p86l.DataFile) {
			df.UsePreRelease = value
		})

		if cacheFile.Releases != nil {
			model.Translate(p86l.ReleasesChangelogText(cacheFile, value), dataFile.Lang)
		}
	})
	if dataFile.UsePreRelease {
		p.prereleaseToggle.SetValue(true)
	} else {
		p.prereleaseToggle.SetValue(false)
	}

	p.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &p.gameVersionText,
			SecondaryWidget: &p.versionText,
		},
		{
			PrimaryWidget:   &p.downloadsText,
			SecondaryWidget: &p.totalDownloadsText,
		},
		{
			PrimaryWidget:   &p.prereleaseText,
			SecondaryWidget: &p.prereleaseToggle,
		},
	})

	linkIcons := [4]*ebiten.Image{assets.IE, assets.Github, assets.Discord, assets.Patreon}
	linkUrls := [4]string{configs.Website, configs.Github, configs.Discord, configs.Patreon}
	for i := range p.linkButtons {
		p.linkButtons[i].SetIcon(linkIcons[i])
		p.linkButtons[i].SetOnDown(func() { model.OpenURL(linkUrls[i]) })
	}

	return nil
}

func (p *Play) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Size: guigui.FixedSize(int(float64(u) * 1.5)),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionHorizontal,
					Items: []guigui.LinearLayoutItem{
						{
							Size: guigui.FlexibleSize(1),
						},
						{
							Widget: &p.linkButtons[0],
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.linkButtons[1],
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.linkButtons[2],
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.linkButtons[3],
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Size: guigui.FlexibleSize(1),
						},
					},
					Gap: u / 2,
				},
			},
			{
				Size: guigui.FixedSize(u),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionHorizontal,
					Items: []guigui.LinearLayoutItem{
						{
							Size: guigui.FlexibleSize(1),
						},
						{
							Widget: &p.actionButtons[0],
							Size:   guigui.FixedSize(u * 4),
						},
						{
							Widget: &p.actionButtons[1],
							Size:   guigui.FixedSize(u * 4),
						},
						{
							Widget: &p.actionButtons[2],
							Size:   guigui.FixedSize(u * 4),
						},
						{
							Size: guigui.FlexibleSize(1),
						},
					},
					Gap: u / 2,
				},
			},
			{
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items: []guigui.LinearLayoutItem{
						{
							Widget: &p.form,
						},
					},
				},
			},
			{
				Widget: &p.changelogPanel,
				Size:   guigui.FlexibleSize(1),
			},
		},
		Gap: u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
