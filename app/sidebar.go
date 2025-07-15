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
	gctx "context"
	"fmt"
	"p86l"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent sidebarContent
}

func (s *Sidebar) SetModel(model *p86l.Model) {
	s.panelContent.SetModel(model)
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetOpacity(&s.panel, 0.9)
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetSize(&s.panelContent, context.ActualSize(s))
	s.panel.SetContent(&s.panelContent)

	appender.AppendChildWidgetWithBounds(&s.panel, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list  basicwidget.List[string]
	stats sidebarStats

	model *p86l.Model
}

func (s *sidebarContent) SetModel(model *p86l.Model) {
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.stats.SetModel(s.model)

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[string]{
		{
			Text: p86l.T("home.title"),
			ID:   "home",
		},
		{
			Text: p86l.T("play.title"),
			ID:   "play",
		},
		{
			Text: p86l.T("changelog.title"),
			ID:   "changelog",
		},
		{
			Text: p86l.T("settings.title"),
			ID:   "settings",
		},
		{
			Text: p86l.T("about.title"),
			ID:   "about",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByID(s.model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			s.model.SetMode("")
			return
		}
		if item.ID == s.model.Mode() {
			return
		}
		s.model.SetMode(item.ID)
	})

	gl := layout.GridLayout{
		Bounds: context.Bounds(s),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
		},
	}
	appender.AppendChildWidgetWithBounds(&s.list, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&s.stats, gl.CellBounds(0, 1))

	return nil
}

func (s *sidebarContent) HandleButtonInput(context *guigui.Context) guigui.HandleInputResult {
	currentIndex := s.list.SelectedItemIndex()
	itemsCount := s.list.ItemsCount()

	if currentIndex >= 0 && currentIndex < itemsCount {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
			newIndex := currentIndex - 1
			if newIndex >= 0 {
				s.list.SelectItemByIndex(newIndex)
				if item, ok := s.list.ItemByIndex(newIndex); ok && item.ID != s.model.Mode() {
					s.model.SetMode(item.ID)
				}
				return guigui.HandleInputByWidget(s)
			}
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
			newIndex := currentIndex + 1
			if newIndex < itemsCount {
				s.list.SelectItemByIndex(newIndex)
				if item, ok := s.list.ItemByIndex(newIndex); ok && item.ID != s.model.Mode() {
					s.model.SetMode(item.ID)
				}
				return guigui.HandleInputByWidget(s)
			}
		}
	}

	return guigui.HandleInputResult{}
}

func (s *sidebarContent) Tick(context *guigui.Context) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		context.SetFocused(s, true)
		return nil
	}
	if context.IsWidgetHitAtCursor(&s.list) {
		_, dy := ebiten.Wheel()

		currentIndex := s.list.SelectedItemIndex()
		itemsCount := s.list.ItemsCount()

		newIndex := currentIndex - int(dy)

		if newIndex < 0 {
			newIndex = 0
		} else if newIndex >= itemsCount {
			newIndex = itemsCount - 1
		}

		if newIndex != currentIndex {
			s.list.SelectItemByIndex(newIndex)
			if item, ok := s.list.ItemByIndex(newIndex); ok && item.ID != s.model.Mode() {
				s.model.SetMode(item.ID)
			}
			context.SetFocused(&s.list, true)
		}
	}

	return nil
}

type sidebarStats struct {
	guigui.DefaultWidget

	progressTextInput basicwidget.TextInput
	toastTextInput    basicwidget.TextInput
	versionText       basicwidget.Text
	ratelimitText     basicwidget.Text

	ratelimitLeft int
	inProgress    bool
	lastTick      int64

	model *p86l.Model
}

func githubRateLimit(ctx gctx.Context) *github.RateLimits {
	limits, _, err := p86l.GithubClient.RateLimit.Get(ctx)
	if err != nil {
		return nil
	}
	return limits
}

func (s *sidebarStats) SetModel(model *p86l.Model) {
	s.model = model
}

func (s *sidebarStats) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.progressTextInput.SetValue(s.model.Progress())
	s.progressTextInput.SetMultiline(true)
	s.progressTextInput.SetAutoWrap(true)
	s.progressTextInput.SetEditable(false)

	if p86l.E.ToastErr != nil {
		s.toastTextInput.SetValue(p86l.E.ToastErr.String())
	} else {
		s.toastTextInput.SetValue("")
	}
	s.toastTextInput.SetMultiline(true)
	s.toastTextInput.SetAutoWrap(true)
	s.toastTextInput.SetEditable(false)

	s.versionText.SetValue(p86l.TheDebugMode.Version)
	s.versionText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	s.versionText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	s.ratelimitText.SetValue(fmt.Sprintf("Ratelimit: %d / 60", s.ratelimitLeft))

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(s),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FlexibleSize(1),
			layout.FixedSize(s.versionText.DefaultSize(context).Y),
			layout.FixedSize(s.ratelimitText.DefaultSize(context).Y),
		},
		Widths: []layout.Size{
			layout.FixedSize(u / 4),
			layout.FlexibleSize(1),
			layout.FixedSize(u / 2),
		},
		RowGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&s.progressTextInput, gl.CellBounds(1, 0))
	appender.AppendChildWidgetWithBounds(&s.toastTextInput, gl.CellBounds(1, 1))
	appender.AppendChildWidgetWithBounds(&s.versionText, gl.CellBounds(1, 2))
	appender.AppendChildWidgetWithBounds(&s.ratelimitText, gl.CellBounds(1, 3))

	return nil
}

func (s *sidebarStats) Tick(context *guigui.Context) error {
	if ebiten.Tick()-s.lastTick >= int64(ebiten.TPS()*2) && !s.inProgress {
		s.lastTick = ebiten.Tick()
		s.inProgress = true

		go func() {
			ctx, cancel := gctx.WithTimeout(gctx.Background(), time.Second*5)
			limits := githubRateLimit(ctx)
			defer cancel()

			if limits != nil {
				s.ratelimitLeft = limits.Core.Remaining
			} else {
				s.ratelimitLeft = -1
			}

			s.inProgress = false
		}()
	}

	return nil
}
