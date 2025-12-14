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
	"sync"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	background basicwidget.Background
	content    sidebarContent
	bottom     bottomContent
}

func (s *Sidebar) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&s.background)
	adder.AddChild(&s.content)
	adder.AddChild(&s.bottom)

	context.SetOpacity(&s.background, 0.9)

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&s.background, widgetBounds.Bounds())
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &s.content,
			},
			{
				Size: guigui.FlexibleSize(1),
			},
			{
				Widget: &s.bottom,
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

type sidebarContent struct {
	guigui.DefaultWidget

	list   basicwidget.List[p86l.SidebarPage]
	bottom bottomContent
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
	layouter.LayoutWidget(&s.list, widgetBounds.Bounds())
}

type bottomContent struct {
	guigui.DefaultWidget

	progressText basicwidget.Text

	rateLimitText basicwidget.Text
	rateLimitStr  string

	sync sync.Once
}

func (b *bottomContent) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&b.progressText)
	adder.AddChild(&b.rateLimitText)

	model := context.Model(b, modelKeyModel).(*p86l.Model)
	b.sync.Do(func() {
		model.SetProgressRefreshFn(func() {
			guigui.RequestRedraw(&b.progressText)
		})
	})

	b.progressText.SetScale(0.8)
	b.progressText.SetAutoWrap(true)
	b.progressText.SetMultiline(true)
	b.progressText.SetValue(model.ProgressText())

	b.rateLimitText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	b.rateLimitText.SetScale(0.8)
	b.rateLimitText.SetAutoWrap(true)
	b.rateLimitText.SetMultiline(true)
	b.rateLimitText.SetValue(b.rateLimitStr)

	return nil
}

func (b *bottomContent) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	model := context.Model(b, modelKeyModel).(*p86l.Model)
	newText := p86l.FormattedCacheExpireText(model.Cache().Get())
	if newText != b.rateLimitStr {
		b.rateLimitStr = newText
		guigui.RequestRedraw(&b.rateLimitText)
	}

	return nil
}

func (b *bottomContent) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &b.progressText,
				Size:   guigui.FlexibleSize(1),
			},
			{
				Widget: &b.rateLimitText,
			},
		},
		Gap: u / 2,
		Padding: guigui.Padding{
			Start: u / 6,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
