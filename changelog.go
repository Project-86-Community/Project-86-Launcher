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
	"time"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/pkg/browser"
)

type Changelog struct {
	guigui.DefaultWidget

	vLayout         internal.VerticalLayout
	changelogText   basicwidget.Text
	changelogButton basicwidget.TextButton

	err error
}

func (c *Changelog) requestChangelog(changelogData *content.Changelog) {
	release, _, err := content.GithubClient.Repositories.GetLatestRelease(content.GithubContext, content.RepoOwner, content.RepoName)
	if err != nil {
		changelogData.Body = err.Error()
		guigui.Disable(&c.changelogButton)
	} else {
		*changelogData = content.Changelog{
			Body:      release.GetBody(),
			URL:       release.GetHTMLURL(),
			Timestamp: time.Now(),
			ExpiresIn: time.Hour,
		}

		cachedJSON, err := json.Marshal(changelogData)
		if err != nil {
			c.err = err
			return
		}
		if err := content.Mgdata.SaveObjectProp("changelog", "changelog.json", cachedJSON); err != nil {
			c.err = err
			return
		}
	}
}

func (c *Changelog) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	if content.Mgdata.ObjectPropExists("changelog", "changelog.json") {
		changelogJSON, err := content.Mgdata.LoadObjectProp("changelog", "changelog.json")
		if err != nil {
			c.err = err
			return
		}
		changelogData := &content.Changelog{}
		err = json.Unmarshal(changelogJSON, &changelogData)
		if err != nil {
			c.err = err
			return
		}

		if content.IsInternet {
			if time.Since(changelogData.Timestamp) > changelogData.ExpiresIn {
				c.requestChangelog(changelogData)
			}
			c.changelogText.SetText(changelogData.Body)
			c.changelogButton.SetText("View changelog")
			guigui.Enable(&c.changelogButton)
			c.changelogButton.SetOnDown(func() {
				go browser.OpenURL(changelogData.URL)
			})
		} else {
			c.changelogText.SetText(changelogData.Body)
			c.changelogButton.SetText("NO INTERNET")
			guigui.Disable(&c.changelogButton)
		}
	} else {
		if content.IsInternet {
			changelogData := content.Changelog{}
			c.requestChangelog(&changelogData)

			c.changelogText.SetText(changelogData.Body)
			c.changelogButton.SetText("View changelog")
			c.changelogButton.SetOnDown(func() {
				go browser.OpenURL(changelogData.URL)
			})
		} else {
			c.changelogText.SetText("NO INTERNET")
			c.changelogButton.SetText("View changelog")
			guigui.Disable(&c.changelogButton)
		}
	}

	u := float64(basicwidget.UnitSize(context))
	w, _ := c.Size(context)
	pt := guigui.Position(c).Add(image.Pt(int(0.5*u), int(0.5*u)))

	c.vLayout.SetHorizontalAlign(internal.HorizontalAlignCenter)

	c.vLayout.SetWidth(context, w-int(1*u))
	guigui.SetPosition(&c.vLayout, pt)

	c.vLayout.SetItems([]*internal.LayoutItem{
		{Widget: &c.changelogText},
		{Widget: &c.changelogButton},
	})
	appender.AppendChildWidget(&c.vLayout)
}

func (c *Changelog) Update(context *guigui.Context) error {
	if c.err != nil {
		fmt.Println(c.err)
		return c.err
	}

	return nil
}

func (c *Changelog) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(c).Size(context)
	w -= sidebarWidth(context)
	return w, h
}
