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
	"p86l"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	panel        basicwidget.Panel
	panelContent sidebarContent
}

func (s *Sidebar) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&s.panel)

	s.panel.SetStyle(basicwidget.PanelStyleSide)
	s.panel.SetBorders(basicwidget.PanelBorder{
		End: true,
	})
	context.SetOpacity(&s.panel, 0.9)
	s.panel.SetContent(&s.panelContent)

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	s.panelContent.setSize(widgetBounds.Bounds().Size())
	layouter.LayoutWidget(&s.panel, widgetBounds.Bounds())
}

type sidebarContent struct {
	guigui.DefaultWidget

	list   basicwidget.List[p86l.SidebarPage]
	bottom sidebarBottom

	size image.Point
}

func (s *sidebarContent) setSize(size image.Point) {
	s.size = size
}

func (s *sidebarContent) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&s.list)
	adder.AddChild(&s.bottom)

	model := context.Model(s, modelKeyModel).(*p86l.Model)
	data := model.Data()
	dataFile := data.Get()

	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.ListItem[p86l.SidebarPage]{
		{
			Text:  p86l.T("home.title"),
			Value: p86l.PageHome,
		},
		{
			Text:  p86l.T("play.play"),
			Value: p86l.PagePlay,
		},
		{
			Text:  p86l.T("settings.title"),
			Value: p86l.PageSettings,
		},
		{
			Text:  p86l.T("about.title"),
			Value: p86l.PageAbout,
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByValue(p86l.SidebarPage(dataFile.Remember.Page))
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			data.Update(func(df *p86l.DataFile) {
				df.Remember.Page = int(p86l.PageHome)
			})
			return
		}
		data.Update(func(df *p86l.DataFile) {
			df.Remember.Page = int(item.Value)
		})
	})

	return nil
}

func (s *sidebarContent) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &s.list,
			},
			{
				Size: guigui.FlexibleSize(1),
			},
			{
				Widget: &s.bottom,
			},
		},
		Gap: u / 2,
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (s *sidebarContent) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return s.size
}

type sidebarBottom struct {
	guigui.DefaultWidget

	progressText    basicwidget.Text
	cacheExpireText basicwidget.Text

	formattedCacheExpireText string
}

func (s *sidebarBottom) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&s.progressText)
	adder.AddChild(&s.cacheExpireText)

	s.progressText.SetAutoWrap(true)

	s.cacheExpireText.SetValue(s.formattedCacheExpireText)
	s.cacheExpireText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	s.cacheExpireText.SetScale(0.8)
	s.cacheExpireText.SetAutoWrap(true)
	s.cacheExpireText.SetMultiline(true)

	return nil
}

func (s *sidebarBottom) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	model := context.Model(s, modelKeyModel).(*p86l.Model)
	newText := p86l.FormattedCacheExpireText(model.Cache().Get())
	if newText != s.formattedCacheExpireText {
		s.formattedCacheExpireText = newText
		guigui.RequestRedraw(s)
	}

	return nil
}

func (s *sidebarBottom) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Size: guigui.FlexibleSize(1),
			},
			{
				Widget: &s.progressText,
			},
			{
				Widget: &s.cacheExpireText,
			},
		},
		Gap: u / 2,
		Padding: guigui.Padding{
			End: u / 6,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
