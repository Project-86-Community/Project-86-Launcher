// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"p86l"
	"p86l/internal/logger"
	"p86l/internal/types"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	background               basicwidget.Background
	versionText              basicwidget.Text
	launchButton             basicwidget.Button
	killButton               basicwidget.Button
	folderButton             basicwidget.Button
	deleteButton             basicwidget.Button
	shortcutButton           basicwidget.Button
	positionSegmentedControl basicwidget.SegmentedControl[types.SidebarPosition]

	layoutItems []guigui.LinearLayoutItem
}

func (s *Sidebar) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&s.background)
	adder.AddWidget(&s.launchButton)
	adder.AddWidget(&s.killButton)
	adder.AddWidget(&s.folderButton)
	adder.AddWidget(&s.deleteButton)
	adder.AddWidget(&s.shortcutButton)
	adder.AddWidget(&s.positionSegmentedControl)

	v, ok := context.Env(s, modelKeyModel)
	if !ok {
		return nil
	}
	model := v.(*p86l.Model)
	t := model.T()
	sidebarModel := model.Sidebar()
	version := sidebarModel.Version()

	context.SetOpacity(&s.background, 0.8)

	hasVersion := version.Tag != ""
	isRunningVersion := sidebarModel.IsRunningVersion(version.Tag, version.OS)

	s.launchButton.SetText(t.Get("home.launch"))
	context.SetEnabled(&s.launchButton, hasVersion && !isRunningVersion && version.Runnable)
	s.launchButton.OnUp(func(context *guigui.Context) {
		sidebarModel.Launch(model.Fake(), s)
	})

	s.killButton.SetText(t.Get("home.kill"))
	context.SetEnabled(&s.killButton, isRunningVersion)
	s.killButton.OnUp(func(context *guigui.Context) {
		sidebarModel.Kill(model.Fake())
	})

	s.folderButton.SetText(t.Get("home.folder"))
	context.SetEnabled(&s.folderButton, hasVersion)
	s.folderButton.OnUp(func(context *guigui.Context) {
		model.Open(version.Tag, false)
	})

	s.deleteButton.SetText(t.Get("home.delete"))
	context.SetEnabled(&s.deleteButton, hasVersion && !isRunningVersion)
	s.deleteButton.OnUp(func(context *guigui.Context) {
		// TODO: handle error properly
		if err := model.DeleteVersion(version.Tag); err != nil {
			logger.Error.Printf("delete failed [%s]: %v", version.Tag, err)
		} else {
			logger.Info.Printf("deleted version: %s", version.Tag)
		}
	})

	s.shortcutButton.SetText(t.Get("home.shortcut"))
	context.SetEnabled(&s.shortcutButton, hasVersion)
	s.shortcutButton.OnUp(func(context *guigui.Context) {
		if err := model.CreateShortcut(version); err != nil {
			logger.Warn.Printf("shortcut failed [%s]: %v", version.Tag, err)
		} else {
			logger.Info.Printf("shortcut created for: %s", version.Tag)
		}
	})

	s.positionSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[types.SidebarPosition]{
		{
			Text:  "◀",
			Value: types.SidebarPositionLeft,
		},
		{
			Text:  "▶",
			Value: types.SidebarPositionRight,
		},
	})
	s.positionSegmentedControl.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := s.positionSegmentedControl.ItemByIndex(index)
		if !ok {
			return
		}
		sidebarModel.SetSidebarPosition(item.Value)
	})
	s.positionSegmentedControl.SelectItemByValue(sidebarModel.SidebarPosition())

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&s.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	s.layoutItems = slices.Delete(s.layoutItems, 0, len(s.layoutItems))
	s.layoutItems = append(s.layoutItems,
		guigui.LinearLayoutItem{
			Widget: &s.launchButton,
		},
		guigui.LinearLayoutItem{
			Widget: &s.killButton,
		},
		guigui.LinearLayoutItem{
			Widget: &s.folderButton,
		},
		guigui.LinearLayoutItem{
			Widget: &s.deleteButton,
		},
		guigui.LinearLayoutItem{
			Widget: &s.shortcutButton,
		},
		guigui.LinearLayoutItem{
			Size: guigui.FlexibleSize(1),
		},
		guigui.LinearLayoutItem{
			Widget: &s.positionSegmentedControl,
		},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     s.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
