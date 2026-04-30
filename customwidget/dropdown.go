// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package customwidget

import (
	"image"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type DropdownItem[T comparable] struct {
	Text     string
	Value    T
	Disabled bool
	Border   bool
}

type Dropdown[T comparable] struct {
	guigui.DefaultWidget

	button    basicwidget.Button
	popupMenu basicwidget.PopupMenu[int]

	items      []DropdownItem[T]
	onSelected func(context *guigui.Context, index int)
	label      string
}

func (d *Dropdown[T]) SetLabel(label string) {
	d.label = label
}

func (d *Dropdown[T]) SetItems(items []DropdownItem[T]) {
	d.items = items
}

func (d *Dropdown[T]) OnItemSelected(fn func(context *guigui.Context, index int)) {
	d.onSelected = fn
}

func (d *Dropdown[T]) ItemByIndex(index int) (DropdownItem[T], bool) {
	if index < 0 || index >= len(d.items) {
		return DropdownItem[T]{}, false
	}
	return d.items[index], true
}

func (d *Dropdown[T]) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&d.button)
	adder.AddWidget(&d.popupMenu)

	d.button.SetText(d.label)
	d.button.OnUp(func(context *guigui.Context) {
		d.popupMenu.SetOpen(true)
	})

	menuItems := make([]basicwidget.PopupMenuItem[int], 0, len(d.items))
	for i, item := range d.items {
		menuItems = append(menuItems, basicwidget.PopupMenuItem[int]{
			Text:     item.Text,
			Value:    i,
			Disabled: item.Disabled,
			Border:   item.Border,
		})
	}
	d.popupMenu.SetItems(menuItems)
	d.popupMenu.OnItemSelected(func(context *guigui.Context, index int) {
		if d.onSelected != nil {
			d.onSelected(context, index)
		}
	})

	return nil
}

func (d *Dropdown[T]) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(&d.button, b)

	u := basicwidget.UnitSize(context)
	menuSize := d.popupMenu.Measure(context, guigui.Constraints{})
	menuPos := image.Pt(b.Min.X, b.Max.Y+u/4)
	layouter.LayoutWidget(&d.popupMenu, image.Rectangle{
		Min: menuPos,
		Max: menuPos.Add(menuSize),
	})
}

func (d *Dropdown[T]) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return d.button.Measure(context, constraints)
}
