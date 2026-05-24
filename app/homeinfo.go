// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

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
