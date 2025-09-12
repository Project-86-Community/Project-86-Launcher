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
	"fmt"
	"image"
	"p86l"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel           basicwidget.Panel
	panelContent    sidebarContent
	cacheExpireText basicwidget.Text

	mainLayout layout.GridLayout
}

func (s *Sidebar) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.panel)
	appender.AppendChildWidget(&s.cacheExpireText)
}

func (s *Sidebar) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorder(basicwidget.PanelBorder{
		End: true,
	})
	context.SetOpacity(&s.panel, 0.9)
	s.panelContent.setSize(context.Bounds(s).Size())
	s.panel.SetContent(&s.panelContent)

	if cacheData := model.Cache().Data(); cacheData != nil && cacheData.RateLimit2 != nil {
		s.cacheExpireText.SetValue(fmt.Sprintf(
			"%d / %d - requests - %s",
			cacheData.RateLimit2.Remaining,
			cacheData.RateLimit2.Limit,
			humanize.RelTime(time.Now(), model.Cache().ExpiresAt(), "remaining", "ago"),
		))
	}
	s.cacheExpireText.SetAutoWrap(true)
	s.cacheExpireText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

	u := basicwidget.UnitSize(context)
	s.mainLayout = layout.GridLayout{
		Bounds: context.Bounds(s),
		Heights: []layout.Size{
			layout.FixedSize(s.panelContent.Measure(context, guigui.FixedWidthConstraints(context.Bounds(s).Dx()-u)).Y),
			layout.FlexibleSize(1),
			layout.FixedSize(u * 2),
		},
	}

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &s.panel:
		return s.mainLayout.CellBounds(0, 0)
	case &s.cacheExpireText:
		return s.mainLayout.CellBounds(0, 2).Inset(basicwidget.UnitSize(context) / 4)
	}
	return image.Rectangle{}
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[string]

	size image.Point
}

func (s *sidebarContent) setSize(size image.Point) {
	s.size = size
}

func (s *sidebarContent) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&s.list)
}

func (s *sidebarContent) Build(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[string]{
		{
			Text:  "Home",
			Value: "home",
		},
		{
			Text:  "Play",
			Value: "play",
		},
		{
			Text:  "Settings",
			Value: "settings",
		},
		{
			Text:  "About",
			Value: "about",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByValue(model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			model.SetMode("")
			return
		}
		model.SetMode(item.Value)
	})

	return nil
}

func (s *sidebarContent) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &s.list:
		return context.Bounds(s)
	}
	return image.Rectangle{}
}

func (s *sidebarContent) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return s.size
}
