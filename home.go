/*
 * SPDX-License-Identifier: GPL-3.0-only
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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

package eightysix

import (
	"eightysix/content"
	"eightysix/internal"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/pkg/browser"
)

type Home struct {
	guigui.DefaultWidget

	gameStatus string
	gamePanel  basicwidget.ScrollablePanel
	gameLayout internal.VerticalLayout

	vLayout       internal.VerticalLayout
	banner        basicwidget.Image
	titleText     basicwidget.Text
	form          basicwidget.Form
	gameInfoText  basicwidget.Text
	gameButton    basicwidget.TextButton
	websiteButton basicwidget.TextButton
	githubButton  basicwidget.TextButton
	discordButton basicwidget.TextButton
	patreonButton basicwidget.TextButton

	err error
}

func (h *Home) requestGame(gameFileData *content.GameFile) {
	release, _, err := content.GithubClient.Repositories.GetLatestRelease(content.GithubContext, content.RepoOwner, content.RepoName)
	if err != nil {
		h.gameInfoText.SetText(err.Error())
		guigui.Disable(&h.gameButton)
	} else {
		if len(release.Assets) == 0 {
			h.gameInfoText.SetText(err.Error())
			guigui.Disable(&h.gameButton)
		} else {
			*gameFileData = content.GameFile{
				Tag:       release.GetTagName(),
				Timestamp: time.Now(),
				ExpiresIn: time.Hour,
			}

			for _, asset := range release.Assets {
				if internal.IsValidGameFile(asset.GetName()) {
					gameFileData.URL = asset.GetBrowserDownloadURL()
				}
			}

			cachedJSON, err := json.Marshal(gameFileData)
			if err != nil {
				h.err = err
				return
			}
			if err := content.Mgdata.SaveObjectProp("game", "game.json", cachedJSON); err != nil {
				h.err = err
				return
			}
		}
	}
}

func (h *Home) requestUpdate(newGameFileData *content.GameFile) {
	release, _, err := content.GithubClient.Repositories.GetLatestRelease(content.GithubContext, content.RepoOwner, content.RepoName)
	if err != nil {
		h.gameInfoText.SetText(err.Error())
		guigui.Disable(&h.gameButton)
	} else {
		if len(release.Assets) == 0 {
			h.gameInfoText.SetText(err.Error())
			guigui.Disable(&h.gameButton)
		} else {
			*newGameFileData = content.GameFile{
				Tag:       release.GetTagName(),
				Timestamp: time.Now(),
				ExpiresIn: time.Hour,
			}

			for _, asset := range release.Assets {
				if internal.IsValidGameFile(asset.GetName()) {
					newGameFileData.URL = asset.GetBrowserDownloadURL()
				}
			}

			oldGameFileJSON, err := content.Mgdata.LoadObjectProp("game", "game.json")
			if err != nil {
				h.err = err
				return
			}
			oldGameFileData := &content.GameFile{}
			err = json.Unmarshal(oldGameFileJSON, &oldGameFileData)
			if err != nil {
				h.err = err
				return
			}

			isNewer, err := internal.CheckNewerVersion(oldGameFileData.Tag, newGameFileData.Tag)
			if err != nil {
				h.err = err
				return
			}
			if isNewer {
				content.UpdateGame = true

				cachedJSON, err := json.Marshal(newGameFileData)
				if err != nil {
					h.err = err
					return
				}
				if err := content.Mgdata.SaveObjectProp("game", "game.json", cachedJSON); err != nil {
					h.err = err
					return
				}
			}
		}
	}
}

func (h *Home) gameInstall() {
	if content.Mgdata.ObjectPropExists("game", "game.json") && content.DownloadStatus == -1 {
		content.DownloadStatus = 0

		gameFileJSON, err := content.Mgdata.LoadObjectProp("game", "game.json")
		if err != nil {
			h.err = err
			return
		}
		gameFileData := &content.GameFile{}
		err = json.Unmarshal(gameFileJSON, &gameFileData)
		if err != nil {
			h.err = err
			return
		}

		folderPath := content.Mgdata.ObjectPropPath("darkmode", "darkmode.data")
		folderPath = internal.TrimDarkModePath(folderPath)

		go func() {
			err = internal.DownloadFile(gameFileData.URL, folderPath+"game.zip")
			if err != nil {
				h.err = err
				return
			}
			internal.ExtractZip(folderPath+"game.zip", folderPath+"run")
			content.DownloadStatus = -1
		}()
	}
}

func (h *Home) gamePlay() {
	folderPath := content.Mgdata.ObjectPropPath("darkmode", "darkmode.data")
	folderPath = internal.TrimDarkModePath(folderPath)

	exePath, err := internal.FindExecutable(folderPath+"run", "Project-86.exe")
	if err != nil {
		h.err = err
		return
	}

	err = internal.RunExecutable(exePath)
	if err != nil {
		h.err = err
		return
	}

	os.Exit(1)
}

func (h *Home) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	img, err := content.TheImageCache.Get("banner", context.ColorMode())
	if err != nil {
		h.err = err
		return
	}
	h.banner.SetImage(img)
	h.gameInfoText.SetText("")

	if content.Mgdata.ObjectPropExists("game", "game.json") {
		gameFileJSON, err := content.Mgdata.LoadObjectProp("game", "game.json")
		if err != nil {
			h.err = err
			return
		}
		gameFileData := &content.GameFile{}
		err = json.Unmarshal(gameFileJSON, &gameFileData)
		if err != nil {
			h.err = err
			return
		}

		folderPath := content.Mgdata.ObjectPropPath("darkmode", "darkmode.data")
		folderPath = internal.TrimDarkModePath(folderPath)

		if internal.FolderExists(folderPath + "run") {
			h.gameButton.SetText("Play")
			guigui.Enable(&h.gameButton)
			h.gameStatus = "play"
		} else {
			h.gameButton.SetText("Install")
			guigui.Enable(&h.gameButton)
			h.gameStatus = "install"
		}

		if content.IsInternet {
			if time.Since(gameFileData.Timestamp) > gameFileData.ExpiresIn {
				h.requestUpdate(gameFileData)
			}
			if content.UpdateGame {
				h.gameButton.SetText("Update")
				guigui.Enable(&h.gameButton)
				h.gameStatus = "update"
			}
			if content.DownloadStatus != -1 {
				if content.DownloadStatus >= 99.99 {
					h.gameInfoText.SetText("Extracting zip file...")
				} else {
					h.gameInfoText.SetText(fmt.Sprintf("%.1f%% - ETA: %s", content.DownloadStatus, content.DownloadETA))
				}
				h.gameButton.SetText("Downloading...")
				guigui.Disable(&h.gameButton)
			}
		}
	} else {
		if content.IsInternet {
			gameFileData := content.GameFile{}
			h.requestGame(&gameFileData)

			h.gameButton.SetText("Install")
			guigui.Disable(&h.gameButton)
		} else {
			h.gameButton.SetText("NO INTERNET")
			guigui.Disable(&h.gameButton)
		}
	}

	u := float64(basicwidget.UnitSize(context))
	w, _ := h.Size(context)
	pt := guigui.Position(h).Add(image.Pt(int(0.5*u), int(0.5*u)))

	{
		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()
		aspectRatio := float64(imgHeight) / float64(imgWidth)

		newWidth := w - int(u)
		newHeight := int(float64(newWidth) * aspectRatio)

		h.banner.SetSize(context, newWidth, newHeight)
	}

	h.titleText.SetBold(true)
	h.titleText.SetScale(2)
	h.titleText.SetText("Welcome to Project 86")
	h.titleText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

	_, gameTextHeight := h.gameInfoText.Size(context)
	h.gamePanel.SetSize(context, w, gameTextHeight+int(2*u))
	h.gameButton.SetWidth(240)

	h.gamePanel.SetContent(func(context *guigui.Context, childAppender *basicwidget.ContainerChildWidgetAppender, offsetX, offsetY float64) {
		p := guigui.Position(&h.gamePanel).Add(image.Pt(int(offsetX), int(offsetY)))

		h.gameLayout.SetHorizontalAlign(internal.HorizontalAlignCenter)
		h.gameLayout.DisableBackground(true)
		h.gameLayout.DisableLineBreak(true)
		h.gameLayout.DisableBorder(true)

		h.gameLayout.SetWidth(context, w-int(1*u))
		guigui.SetPosition(&h.gameLayout, image.Pt(p.X+int(u), p.Y+int(u)))

		h.gameLayout.SetItems([]*internal.LayoutItem{
			{Widget: &h.gameInfoText},
		})
		childAppender.AppendChildWidget(&h.gameLayout)
	})
	h.gamePanel.SetPadding(int(2*u), 0)

	h.websiteButton.SetText("Website")
	h.websiteButton.SetWidth(110)

	h.githubButton.SetText("Github")
	h.githubButton.SetWidth(110)

	h.discordButton.SetText("Discord")
	h.discordButton.SetWidth(110)

	h.patreonButton.SetText("Patreon")
	h.patreonButton.SetWidth(110)

	titleTextWidth, _ := h.titleText.Size(context)
	h.form.SetWidth(context, titleTextWidth/2)
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

	h.vLayout.SetHorizontalAlign(internal.HorizontalAlignCenter)
	h.vLayout.DisableBackground(true)
	h.vLayout.DisableLineBreak(true)
	h.vLayout.DisableBorder(true)

	h.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&h.vLayout, pt)

	h.vLayout.SetItems([]*internal.LayoutItem{
		{Widget: &h.banner},
		{Widget: &h.titleText},
		{Widget: &h.gamePanel},
		{Widget: &h.gameButton},
		{Widget: &h.form},
	})
	appender.AppendChildWidget(&h.vLayout)

	h.gameButton.SetOnDown(func() {
		switch h.gameStatus {
		case "install":
			h.gameInstall()
		case "update":
			content.UpdateGame = false
			gameFileData := content.GameFile{}
			h.requestGame(&gameFileData)
			h.gameInstall()
		case "play":
			h.gamePlay()
		}
	})
	h.websiteButton.SetOnDown(func() {
		go browser.OpenURL("https://taliayaya.github.io/Project-86-Website/")
	})
	h.githubButton.SetOnDown(func() {
		go browser.OpenURL("https://github.com/Taliayaya/Project-86")
	})
	h.discordButton.SetOnDown(func() {
		go browser.OpenURL("https://discord.gg/Yh2TQH97yA")
	})
	h.patreonButton.SetOnDown(func() {
		go browser.OpenURL("https://patreon.com/project86")
	})
}

func (h *Home) Update(context *guigui.Context) error {
	if h.err != nil {
		return h.err
	}
	return nil
}

func (h *Home) Size(context *guigui.Context) (int, int) {
	w, _h := guigui.Parent(h).Size(context)
	w -= sidebarWidth(context)
	return w, _h
}
