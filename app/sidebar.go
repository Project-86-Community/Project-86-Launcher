// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"p86l/internal/logger"
	"p86l/internal/service"
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

	launch, ok1 := envMust[service.LaunchService](context, s, keyLaunch)
	fileOpen, ok2 := envMust[service.FileOpenerService](context, s, keyFileOpen)
	shortcutService, ok3 := envMust[service.ShortcutService](context, s, keyShortcut)
	versionService, ok4 := envMust[service.VersionService](context, s, keyVersion)
	model, ok5 := envMust[*Model](context, s, modelKeyModel)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return nil
	}
	t := model.T()

	ver := launch.CurrentVersion()

	context.SetOpacity(&s.background, 0.8)

	hasVersion := ver.Tag != ""
	isRunningVersion := launch.IsRunningVersion(ver.Tag, ver.OS)

	s.launchButton.SetText(t.Get("home.launch"))
	context.SetEnabled(&s.launchButton, hasVersion && !isRunningVersion && ver.Runnable)
	s.launchButton.OnUp(func(context *guigui.Context) {
		launch.Launch(ver, func() {
			guigui.RequestRebuild(s)
		})
	})

	s.killButton.SetText(t.Get("home.kill"))
	context.SetEnabled(&s.killButton, isRunningVersion)
	s.killButton.OnUp(func(context *guigui.Context) {
		launch.Kill()
	})

	s.folderButton.SetText(t.Get("home.folder"))
	context.SetEnabled(&s.folderButton, hasVersion)
	s.folderButton.OnUp(func(context *guigui.Context) {
		fileOpen.OpenPath(ver.Tag)
	})

	s.deleteButton.SetText(t.Get("home.delete"))
	context.SetEnabled(&s.deleteButton, hasVersion && !isRunningVersion)
	s.deleteButton.OnUp(func(context *guigui.Context) {
		if err := versionService.DeleteVersion(ver.Tag); err != nil {
			logger.Error.Printf("delete failed [%s]: %v", ver.Tag, err)
		} else {
			logger.Info.Printf("deleted version: %s", ver.Tag)
		}
	})

	s.shortcutButton.SetText(t.Get("home.shortcut"))
	context.SetEnabled(&s.shortcutButton, hasVersion)
	s.shortcutButton.OnUp(func(context *guigui.Context) {
		if err := shortcutService.CreateShortcut(ver); err != nil {
			logger.Warn.Printf("shortcut failed [%s]: %v", ver.Tag, err)
		} else {
			logger.Info.Printf("shortcut created for: %s", ver.Tag)
		}
	})

	s.positionSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[types.SidebarPosition]{
		{Text: "◀", Value: types.SidebarPositionLeft},
		{Text: "▶", Value: types.SidebarPositionRight},
	})
	s.positionSegmentedControl.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := s.positionSegmentedControl.ItemByIndex(index)
		if !ok {
			return
		}
		model.SetSidebarPosition(item.Value)
	})
	s.positionSegmentedControl.SelectItemByValue(model.SidebarPosition())

	return nil
}

func (s *Sidebar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&s.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	s.layoutItems = slices.Delete(s.layoutItems, 0, len(s.layoutItems))
	s.layoutItems = append(s.layoutItems,
		guigui.LinearLayoutItem{Widget: &s.launchButton},
		guigui.LinearLayoutItem{Widget: &s.killButton},
		guigui.LinearLayoutItem{Widget: &s.folderButton},
		guigui.LinearLayoutItem{Widget: &s.deleteButton},
		guigui.LinearLayoutItem{Widget: &s.shortcutButton},
		guigui.LinearLayoutItem{Size: guigui.FlexibleSize(1)},
		guigui.LinearLayoutItem{Widget: &s.positionSegmentedControl},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     s.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start: u / 2, Top: u / 2, End: u / 2, Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
