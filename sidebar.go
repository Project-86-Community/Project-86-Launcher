/*
 * SPDX-License-Identifier: GPL-3.0-only
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
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	sidebar        basicwidget.Sidebar
	sidebarContent sidebarContent
}

func (s *Sidebar) SetModel(model *Model) {
	s.sidebarContent.SetModel(model)
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetSize(&s.sidebarContent, context.Size(s))
	s.sidebar.SetContent(&s.sidebarContent)

	appender.AppendChildWidgetWithBounds(&s.sidebar, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.TextList[string]

	model *Model
}

func (s *sidebarContent) SetModel(model *Model) {
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	items := []basicwidget.TextListItem[string]{
		{
			Text: l.T("home.title"),
			Tag:  "home",
		},
		{
			Text: l.T("changelog.title"),
			Tag:  "changelog",
		},
		{
			Text: l.T("settings.title"),
			Tag:  "settings",
		},
		{
			Text: l.T("about.title"),
			Tag:  "about",
		},
		{
			Text: l.T("back.title"),
			Tag:  "back",
		},
	}

	s.list.SetItems(items)
	s.list.SelectItemByTag(s.model.Mode())
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	s.list.SetOnItemSelected(func(index int) {
		item, ok := s.list.ItemByIndex(index)
		if !ok {
			s.model.SetMode("")
			return
		}
		s.model.SetMode(item.Tag)
	})
	appender.AppendChildWidgetWithBounds(&s.list, context.Bounds(s))

	return nil
}
