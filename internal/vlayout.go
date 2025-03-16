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

package internal

import (
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type HorizontalAlign int

const (
	HorizontalAlignStart HorizontalAlign = iota
	HorizontalAlignCenter
	HorizontalAlignEnd
)

type LayoutItem struct {
	Widget guigui.Widget
}

type VerticalLayout struct {
	guigui.DefaultWidget

	hAlign     HorizontalAlign
	background bool
	lineBreak  bool
	border     bool

	items             []*LayoutItem
	widthMinusDefault int

	widgetBounds []image.Rectangle
}

func formItemPadding(context *guigui.Context) (int, int) {
	return basicwidget.UnitSize(context) / 2, basicwidget.UnitSize(context) / 4
}

func (v *VerticalLayout) SetItems(items []*LayoutItem) {
	v.items = slices.Delete(v.items, 0, len(v.items))
	v.items = append(v.items, items...)
}

func (v *VerticalLayout) SetHorizontalAlign(align HorizontalAlign) {
	if v.hAlign == align {
		return
	}

	v.hAlign = align
	guigui.RequestRedraw(v)
}

func (v *VerticalLayout) DisableBackground(value bool) {
	if v.background == value {
		return
	}

	v.background = value
	guigui.RequestRedraw(v)
}

func (v *VerticalLayout) DisableLineBreak(value bool) {
	if v.lineBreak == value {
		return
	}

	v.lineBreak = value
	guigui.RequestRedraw(v)
}

func (v *VerticalLayout) DisableBorder(value bool) {
	if v.border == value {
		return
	}

	v.border = value
	guigui.RequestRedraw(v)
}

func (v *VerticalLayout) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	v.calcItemBounds(context)

	for i, item := range v.items {
		v.widgetBounds = append(v.widgetBounds, image.Rectangle{})

		if item.Widget == nil {
			continue
		}

		if item.Widget != nil {
			switch v.hAlign {
			case HorizontalAlignStart:
				guigui.SetPosition(item.Widget, v.widgetBounds[i].Min)
			case HorizontalAlignCenter:
				u := float64(basicwidget.UnitSize(context))
				bounds := v.widgetBounds[i]
				w, _ := v.Size(context)
				iw, _ := item.Widget.Size(context)
				centerX := bounds.Min.X + ((w / 2) - (iw / 2)) - int(0.5*u)
				guigui.SetPosition(item.Widget, image.Point{X: centerX, Y: bounds.Min.Y})
			case HorizontalAlignEnd:
				u := float64(basicwidget.UnitSize(context))
				bounds := v.widgetBounds[i]
				w, _ := v.Size(context)
				iw, _ := item.Widget.Size(context)
				centerX := bounds.Min.X + (w - iw) - int(0.5*u)*2
				guigui.SetPosition(item.Widget, image.Point{X: centerX, Y: bounds.Min.Y})
			default:
				guigui.SetPosition(item.Widget, v.widgetBounds[i].Min)
			}
			appender.AppendChildWidget(item.Widget)
		}
	}
}

func (v *VerticalLayout) calcItemBounds(context *guigui.Context) {
	v.widgetBounds = slices.Delete(v.widgetBounds, 0, len(v.widgetBounds))

	paddingX, paddingY := formItemPadding(context)

	var y int
	for i, item := range v.items {
		v.widgetBounds = append(v.widgetBounds, image.Rectangle{})

		if item.Widget == nil {
			continue
		}

		var widgetH int
		if item.Widget != nil {
			_, widgetH = item.Widget.Size(context)
		}
		h := max(widgetH, minFormItemHeight(context))
		baseBounds := guigui.Bounds(v)
		baseBounds.Min.X += paddingX
		baseBounds.Max.X -= paddingX
		baseBounds.Min.Y += y
		baseBounds.Max.Y = baseBounds.Min.Y + h

		if item.Widget != nil {
			bounds := baseBounds
			ww, wh := item.Widget.Size(context)
			bounds.Max.X = bounds.Min.X + ww
			pY := (h + 2*paddingY - wh) / 2
			if wh < basicwidget.UnitSize(context)+2*paddingY {
				pY = min(pY, max(0, (basicwidget.UnitSize(context)+2*paddingY-wh)/2))
			}
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			v.widgetBounds[i] = bounds
		}

		y += h + 2*paddingY
	}
}

func (v *VerticalLayout) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := guigui.Bounds(v)
	bounds.Max.Y = bounds.Min.Y + v.height(context)
	if !v.background {
		basicwidget.DrawRoundedRect(context, dst, bounds, basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.925), basicwidget.RoundedCornerRadius(context))
	}

	if !v.lineBreak && len(v.items) > 0 {
		paddingX, paddingY := formItemPadding(context)
		y := paddingY
		for _, item := range v.items[:len(v.items)-1] {
			var widgetH int
			if item.Widget != nil {
				_, widgetH = item.Widget.Size(context)
			}
			h := max(widgetH, minFormItemHeight(context))
			y += h + 2*paddingY

			x0 := float32(bounds.Min.X + paddingX)
			x1 := float32(bounds.Max.X - paddingX)
			y := float32(y) + float32(paddingY)
			width := 1 * float32(context.Scale())
			clr := basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.875)
			vector.StrokeLine(dst, x0, y, x1, y, width, clr, false)
		}
	}

	if !v.border {
		basicwidget.DrawRoundedRectBorder(context, dst, bounds, basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.875), basicwidget.RoundedCornerRadius(context), 1*float32(context.Scale()), basicwidget.RoundedRectBorderTypeRegular)
	}
}

func (v *VerticalLayout) SetWidth(context *guigui.Context, width int) {
	v.widthMinusDefault = width - defaultFormWidth(context)
}

func (v *VerticalLayout) Size(context *guigui.Context) (int, int) {
	return v.widthMinusDefault + defaultFormWidth(context), v.height(context)
}

func defaultFormWidth(context *guigui.Context) int {
	return 6 * basicwidget.UnitSize(context)
}

func (v *VerticalLayout) height(context *guigui.Context) int {
	_, paddingY := formItemPadding(context)

	var y int
	for _, item := range v.items {
		if item.Widget == nil || !guigui.IsVisible(item.Widget) {
			continue
		}
		var widgetH int
		if item.Widget != nil {
			_, widgetH = item.Widget.Size(context)
		}
		h := max(widgetH, minFormItemHeight(context))
		y += h + 2*paddingY
	}
	return y
}

func minFormItemHeight(context *guigui.Context) int {
	return basicwidget.UnitSize(context)
}
