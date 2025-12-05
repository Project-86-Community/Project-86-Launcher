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

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Play struct {
	guigui.DefaultWidget

	installButton, updateButton, playButton                         basicwidget.Button
	form                                                            basicwidget.Form
	gameVersionText, versionText, downloadsText, totalDownloadsText basicwidget.Text
	prereleaseText                                                  basicwidget.Text
	prereleaseToggle                                                basicwidget.Toggle
	changelogPanel                                                  basicwidget.Panel
	changelogText                                                   basicwidget.Text
	websiteButton, githubButton, discordButton, patreonButton       basicwidget.Button
}

func (p *Play) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&p.installButton)
	adder.AddChild(&p.updateButton)
	adder.AddChild(&p.playButton)
	adder.AddChild(&p.form)
	adder.AddChild(&p.changelogPanel)
	adder.AddChild(&p.websiteButton)
	adder.AddChild(&p.githubButton)
	adder.AddChild(&p.discordButton)
	adder.AddChild(&p.patreonButton)

	model := context.Model(p, modelKeyModel).(*p86l.Model)
	data := model.Data()
	dataFile := data.Get()
	cacheFile := model.Cache().Get()

	p.installButton.SetText(p86l.T("play.install"))
	p.updateButton.SetText(p86l.T("play.update"))
	p.playButton.SetText(p86l.T("play.play"))

	if dataFile.UsePreRelease {
		preReleaseAvail, _ := model.CheckFilesCached(p86l.PathGamePreRelease)
		context.SetEnabled(&p.playButton, preReleaseAvail)
	} else {
		stableAvail, _ := model.CheckFilesCached(p86l.PathGameStable)
		context.SetEnabled(&p.playButton, stableAvail)
	}

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

		if cacheFile.Releases != nil && dataFile.TranslateChangelog && dataFile.Lang != "en" {
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

	p.websiteButton.SetOnDown(func() { model.OpenURL(configs.Website) })
	p.githubButton.SetOnDown(func() { model.OpenURL(configs.Github) })
	p.discordButton.SetOnDown(func() { model.OpenURL(configs.Discord) })
	p.patreonButton.SetOnDown(func() { model.OpenURL(configs.Patreon) })

	p.websiteButton.SetIcon(assets.IE)
	p.githubButton.SetIcon(assets.Github)
	p.discordButton.SetIcon(assets.Discord)
	p.patreonButton.SetIcon(assets.Patreon)

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
							Widget: &p.websiteButton,
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.githubButton,
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.discordButton,
							Size:   guigui.FixedSize(int(float64(u) * 1.5)),
						},
						{
							Widget: &p.patreonButton,
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
							Widget: &p.installButton,
							Size:   guigui.FixedSize(u * 4),
						},
						{
							Widget: &p.updateButton,
							Size:   guigui.FixedSize(u * 4),
						},
						{
							Widget: &p.playButton,
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
