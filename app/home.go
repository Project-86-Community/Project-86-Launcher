// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"fmt"
	"p86l/internal/service"
	"p86l/internal/types"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Home struct {
	guigui.DefaultWidget

	sidebar                  Sidebar
	positionSegmentedControl basicwidget.SegmentedControl[types.ListPosition]
	list                     homeList
	info                     homeInfo
	noVersions               basicwidget.Text

	contentItems []guigui.LinearLayoutItem
	layoutItems  []guigui.LinearLayoutItem
}

func (h *Home) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&h.sidebar)
	adder.AddWidget(&h.positionSegmentedControl)
	adder.AddWidget(&h.list)
	adder.AddWidget(&h.info)
	adder.AddWidget(&h.noVersions)

	launchService, ok1 := envMust[service.LaunchService](context, h, keyLaunch)
	model, ok2 := envMust[*Model](context, h, modelKeyModel)
	if !ok1 || !ok2 {
		return nil
	}

	h.positionSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[types.ListPosition]{
		{Text: "▲", Value: types.ListPositionTop},
		{Text: "▼", Value: types.ListPositionBottom},
	})
	h.positionSegmentedControl.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := h.positionSegmentedControl.ItemByIndex(index)
		if !ok {
			return
		}
		model.SetListPosition(item.Value)
	})
	h.positionSegmentedControl.SelectItemByValue(model.ListPosition())

	h.list.onVersionSelected = func(ver service.Version) {
		launchService.SetCurrentVersion(ver)
	}

	return nil
}

func (h *Home) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	model, ok := envMust[*Model](context, h, modelKeyModel)
	if !ok {
		return
	}

	u := basicwidget.UnitSize(context)
	sidebarItem := guigui.LinearLayoutItem{
		Widget: &h.sidebar,
		Size:   guigui.FixedSize(u * 8),
	}

	h.contentItems = slices.Delete(h.contentItems, 0, len(h.contentItems))
	h.contentItems = append(h.contentItems,
		guigui.LinearLayoutItem{Widget: &h.positionSegmentedControl},
	)
	if model.ListPosition() == types.ListPositionTop {
		h.contentItems = append(h.contentItems,
			guigui.LinearLayoutItem{Widget: &h.list, Size: guigui.FlexibleSize(1)},
			guigui.LinearLayoutItem{Widget: &h.info, Size: guigui.FlexibleSize(2)},
		)
	} else {
		h.contentItems = append(h.contentItems,
			guigui.LinearLayoutItem{Widget: &h.info, Size: guigui.FlexibleSize(2)},
			guigui.LinearLayoutItem{Widget: &h.list, Size: guigui.FlexibleSize(1)},
		)
	}
	contentItem := guigui.LinearLayoutItem{
		Size: guigui.FlexibleSize(1),
		Layout: guigui.LinearLayout{
			Direction: guigui.LayoutDirectionVertical,
			Items:     h.contentItems,
			Gap:       u / 2,
			Padding: guigui.Padding{
				Start: u / 2, Top: u / 2, End: u / 2, Bottom: u / 2,
			},
		},
	}

	h.layoutItems = slices.Delete(h.layoutItems, 0, len(h.layoutItems))
	if model.SidebarPosition() == types.SidebarPositionLeft {
		h.layoutItems = append(h.layoutItems, sidebarItem, contentItem)
	} else {
		h.layoutItems = append(h.layoutItems, contentItem, sidebarItem)
	}
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     h.layoutItems,
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

type homeList struct {
	guigui.DefaultWidget

	background basicwidget.Background
	titleText  basicwidget.Text
	list       basicwidget.List[string]
	noVersions basicwidget.Text

	versions    []service.Version
	items       []basicwidget.ListItem[string]
	layoutItems []guigui.LinearLayoutItem

	onVersionSelected func(ver service.Version)
}

func (h *homeList) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&h.background)
	adder.AddWidget(&h.titleText)
	adder.AddWidget(&h.list)
	if len(h.versions) == 0 {
		adder.AddWidget(&h.noVersions)
	}

	downloadService, ok1 := envMust[service.DownloadService](context, h, keyDownload)
	model, ok2 := envMust[*Model](context, h, modelKeyModel)
	if !ok1 || !ok2 {
		return nil
	}
	t := model.T()

	versions, err := downloadService.InstalledVersions()
	if err != nil {
		return fmt.Errorf("home: failed to list versions: %w", err)
	}
	h.versions = versions

	context.SetOpacity(&h.background, 0.8)

	h.titleText.SetValue(t.Get("home.title"))
	h.titleText.SetScale(1.2)

	h.items = slices.Delete(h.items, 0, len(h.items))
	for i, ver := range h.versions {
		text := ver.Tag
		osLabel := ""
		switch ver.OS {
		case "windows":
			osLabel = "Windows"
		case "darwin":
			osLabel = "macOS"
		case "linux":
			osLabel = "Linux"
		}
		if osLabel != "" {
			text += " (" + osLabel + ")"
		}
		if !ver.Runnable {
			text += " - " + t.Get("home.incompatible")
		}
		// Use index as value to ensure uniqueness when the same tag has
		// multiple OS builds (e.g. v1.0.0 for both Windows and Linux).
		h.items = append(h.items, basicwidget.ListItem[string]{
			Text: text, Value: fmt.Sprintf("%s:%d", ver.Tag, i),
		})
	}
	h.list.SetItems(h.items)
	h.list.SetItemHeight(basicwidget.UnitSize(context))
	h.list.OnItemSelected(func(context *guigui.Context, index int) {
		if index < 0 || index >= len(h.versions) {
			return
		}
		if h.onVersionSelected != nil {
			h.onVersionSelected(h.versions[index])
		}
	})

	h.noVersions.SetValue(t.Get("home.no_versions"))
	h.noVersions.SetWrapMode(basicwidget.WrapModeNormal)

	return nil
}

func (h *homeList) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	h.layoutItems = slices.Delete(h.layoutItems, 0, len(h.layoutItems))
	if len(h.versions) == 0 {
		h.layoutItems = append(h.layoutItems,
			guigui.LinearLayoutItem{Widget: &h.noVersions, Size: guigui.FlexibleSize(1)},
		)
	} else {
		h.layoutItems = append(h.layoutItems,
			guigui.LinearLayoutItem{Widget: &h.titleText},
			guigui.LinearLayoutItem{Widget: &h.list, Size: guigui.FlexibleSize(1)},
		)
	}
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &h.background,
				Size:   guigui.FlexibleSize(1),
				Layout: guigui.LinearLayout{
					Direction: guigui.LayoutDirectionVertical,
					Items:     h.layoutItems,
					Gap:       u / 8,
					Padding: guigui.Padding{
						Start: u / 4, Top: u / 4, End: u / 4, Bottom: u / 4,
					},
				},
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

type homeInfo struct {
	guigui.DefaultWidget

	background basicwidget.Background
	text       basicwidget.Text

	layoutItems []guigui.LinearLayoutItem
}

func (h *homeInfo) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&h.background)
	adder.AddWidget(&h.text)

	model, ok := envMust[*Model](context, h, modelKeyModel)
	if !ok {
		return nil
	}
	t := model.T()

	context.SetOpacity(&h.background, 0.8)
	h.text.SetValue(t.Get("home.info"))
	h.text.SetWrapMode(basicwidget.WrapModeNormal)

	return nil
}

func (h *homeInfo) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	h.layoutItems = slices.Delete(h.layoutItems, 0, len(h.layoutItems))
	h.layoutItems = append(h.layoutItems,
		guigui.LinearLayoutItem{
			Widget: &h.background,
			Size:   guigui.FlexibleSize(1),
			Layout: guigui.LinearLayout{
				Direction: guigui.LayoutDirectionVertical,
				Items: []guigui.LinearLayoutItem{
					{Widget: &h.text},
				},
			},
		},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     h.layoutItems,
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
