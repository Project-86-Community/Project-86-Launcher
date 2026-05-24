// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type BottomBar struct {
	guigui.DefaultWidget

	background basicwidget.Background
	text1      basicwidget.Text
	text2      basicwidget.Text

	layoutItems []guigui.LinearLayoutItem
}

func (b *BottomBar) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&b.background)
	adder.AddWidget(&b.text1)
	adder.AddWidget(&b.text2)

	b.text1.SetValue("Hello World")
	b.text2.SetValue("Hello World")

	return nil
}

func (b *BottomBar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&b.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	b.layoutItems = slices.Delete(b.layoutItems, 0, len(b.layoutItems))
	b.layoutItems = append(b.layoutItems,
		guigui.LinearLayoutItem{
			Widget: &b.text1,
		},
		guigui.LinearLayoutItem{
			Widget: &b.text2,
		},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     b.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start:  u / 8,
			Top:    u / 8,
			End:    u / 8,
			Bottom: u / 8,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
