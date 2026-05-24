// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
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
