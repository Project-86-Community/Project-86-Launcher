// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"fmt"
	"p86l/internal/service"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

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
