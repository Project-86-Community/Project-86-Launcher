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

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent sidebarContent
	bottom       sidebarBottom
}

func (s *Sidebar) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&s.panel)
	adder.AddChild(&s.bottom)
}

func (s *Sidebar) Update(context *guigui.Context) error {
	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorders(basicwidget.PanelBorder{
		End: true,
	})
	context.SetOpacity(&s.panel, 0.9)
	s.panelContent.setSize(context.Bounds(s).Size())
	s.panel.SetContent(&s.panelContent)

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	switch widget {
	case &s.panel:
		return context.Bounds(s)
	case &s.bottom:
		u := basicwidget.UnitSize(context)
		return (guigui.LinearLayout{
			Direction: guigui.LayoutDirectionVertical,
			Items: []guigui.LinearLayoutItem{
				{
					Size: guigui.FlexibleSize(1),
				},
				{
					Widget: &s.bottom,
					Size:   guigui.FixedSize(2 * u),
				},
			},
		}).WidgetBounds(context, context.Bounds(s), widget)
	}
	return image.Rectangle{}
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.List[p86l.Pages]

	size image.Point
}

func (s *sidebarContent) setSize(size image.Point) {
	s.size = size
}

func (s *sidebarContent) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&s.list)
}

func (s *sidebarContent) Update(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[p86l.Pages]{
		{
			Text:  "Home",
			Value: p86l.PageHome,
		},
		{
			Text:  "Play",
			Value: p86l.PagePlay,
		},
		{
			Text:  "Settings",
			Value: p86l.PageSettings,
		},
		{
			Text:  "About",
			Value: p86l.PageAbout,
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByValue(model.Data().Page())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			model.Data().SetPage(p86l.PageHome)
			return
		}
		model.Data().SetPage(item.Value)
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

type sidebarBottom struct {
	guigui.DefaultWidget

	cacheExpireText basicwidget.Text
}

func (s *sidebarBottom) AddChildren(context *guigui.Context, adder *guigui.ChildAdder) {
	adder.AddChild(&s.cacheExpireText)
}

func (s *sidebarBottom) Update(context *guigui.Context) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)

	s.cacheExpireText.SetValue(model.Cache().ExpireTimeFormatted())
	s.cacheExpireText.SetAutoWrap(true)
	s.cacheExpireText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

	return nil
}

func (s *sidebarBottom) Layout(context *guigui.Context, widget guigui.Widget) image.Rectangle {
	u := basicwidget.UnitSize(context)
	switch widget {
	case &s.cacheExpireText:
		return context.Bounds(s).Inset(u / 4)
	}
	return image.Rectangle{}
}
