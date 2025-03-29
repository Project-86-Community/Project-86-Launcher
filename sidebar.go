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

package p86l

import (
	"image"
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	sidebar         basicwidget.Sidebar
	list            basicwidget.List
	listItemWidgets []basicwidget.ListItem

	initOnce sync.Once
}

func sidebarWidth(context *guigui.Context) int {
	return 8 * basicwidget.UnitSize(context)
}

func (s *Sidebar) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	_, h := s.Size(context)
	s.sidebar.SetSize(context, sidebarWidth(context), h)
	s.sidebar.SetContent(context, func(context *guigui.Context, childAppender *basicwidget.ContainerChildWidgetAppender, offsetX, offsetY float64) {
		s.list.SetWidth(sidebarWidth(context))
		s.list.SetHeight(h)
		guigui.SetPosition(&s.list, guigui.Position(s).Add(image.Pt(int(offsetX), int(offsetY))))
		childAppender.AppendChildWidget(&s.list)
	})
	guigui.SetPosition(&s.sidebar, guigui.Position(s))
	appender.AppendChildWidget(&s.sidebar)

	s.list.SetStyle(basicwidget.ListStyleSidebar)
	if len(s.listItemWidgets) == 0 {
		{
			var t basicwidget.Text
			t.SetScale(1.2)
			t.SetText("Home")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetWidth(sidebarWidth(context))
			t.SetHeight((basicwidget.UnitSize(context) / 2) + basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "home",
			})
		}
		{
			var t basicwidget.Text
			t.SetScale(1.2)
			t.SetText("Settings")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetWidth(sidebarWidth(context))
			t.SetHeight((basicwidget.UnitSize(context) / 2) + basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "settings",
			})
		}
		{
			var t basicwidget.Text
			t.SetScale(1.2)
			t.SetText("Instances")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetWidth(sidebarWidth(context))
			t.SetHeight((basicwidget.UnitSize(context) / 2) + basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "instances",
			})
		}
		{
			var t basicwidget.Text
			t.SetScale(1.2)
			t.SetText("Changelog")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetWidth(sidebarWidth(context))
			t.SetHeight((basicwidget.UnitSize(context) / 2) + basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "changelog",
			})
		}
		{
			var t basicwidget.Text
			t.SetScale(1.2)
			t.SetText("About")
			t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
			t.SetWidth(sidebarWidth(context))
			t.SetHeight((basicwidget.UnitSize(context) / 2) + basicwidget.UnitSize(context))
			s.listItemWidgets = append(s.listItemWidgets, basicwidget.ListItem{
				Content:    &t,
				Selectable: true,
				Tag:        "about",
			})
		}
	}
	s.list.SetItems(s.listItemWidgets)

	s.initOnce.Do(func() {
		s.list.SetSelectedItemIndex(0)
	})
}

func (s *Sidebar) Update(context *guigui.Context) error {
	for i, w := range s.listItemWidgets {
		t := w.Content.(*basicwidget.Text)
		if s.list.SelectedItemIndex() == i {
			t.SetColor(basicwidget.DefaultActiveListItemTextColor(context))
		} else {
			t.SetColor(basicwidget.DefaultTextColor(context))
		}
	}
	return nil
}

func (s *Sidebar) Size(context *guigui.Context) (int, int) {
	_, h := guigui.Parent(s).Size(context)
	return sidebarWidth(context), h
}

func (s *Sidebar) SelectedItemTag() string {
	item, ok := s.list.SelectedItem()
	if !ok {
		return ""
	}
	return item.Tag.(string)
}

func (s *Sidebar) SetSelectedItemIndex(index int) {
	s.list.SetSelectedItemIndex(index)
}
