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

type VerticalAlign int

const (
	VerticalAlignStart VerticalAlign = iota
	VerticalAlignCenter
	VerticalAlignEnd
)

type HorizontalLayout struct {
	guigui.DefaultWidget

	vAlign     VerticalAlign
	background bool
	lineBreak  bool
	border     bool

	items              []*LayoutItem
	heightMinusDefault int

	widgetBounds []image.Rectangle
}

func (h *HorizontalLayout) SetItems(items []*LayoutItem) {
	h.items = slices.Delete(h.items, 0, len(h.items))
	h.items = append(h.items, items...)
}

func (h *HorizontalLayout) SetVerticalAlign(align VerticalAlign) {
	if h.vAlign == align {
		return
	}

	h.vAlign = align
	guigui.RequestRedraw(h)
}

func (h *HorizontalLayout) DisableBackground(value bool) {
	if h.background == value {
		return
	}

	h.background = value
	guigui.RequestRedraw(h)
}

func (h *HorizontalLayout) DisableLineBreak(value bool) {
	if h.lineBreak == value {
		return
	}

	h.lineBreak = value
	guigui.RequestRedraw(h)
}

func (h *HorizontalLayout) DisableBorder(value bool) {
	if h.border == value {
		return
	}

	h.border = value
	guigui.RequestRedraw(h)
}

func (h *HorizontalLayout) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	h.calcItemBounds(context)

	for i, item := range h.items {
		h.widgetBounds = append(h.widgetBounds, image.Rectangle{})

		if item.Widget == nil {
			continue
		}

		if item.Widget != nil {
			switch h.vAlign {
			case VerticalAlignStart:
				guigui.SetPosition(item.Widget, h.widgetBounds[i].Min)
			case VerticalAlignCenter:
				u := float64(basicwidget.UnitSize(context))
				bounds := h.widgetBounds[i]
				_, h := h.Size(context)
				_, ih := item.Widget.Size(context)
				centerY := bounds.Min.Y + ((h / 2) - (ih / 2)) - int(0.5*u)
				guigui.SetPosition(item.Widget, image.Point{X: bounds.Min.X, Y: centerY})
			case VerticalAlignEnd:
				u := float64(basicwidget.UnitSize(context))
				bounds := h.widgetBounds[i]
				_, h := h.Size(context)
				_, ih := item.Widget.Size(context)
				centerY := bounds.Min.Y + (h - ih) - int(0.5*u)*2
				guigui.SetPosition(item.Widget, image.Point{X: bounds.Min.X, Y: centerY})
			default:
				guigui.SetPosition(item.Widget, h.widgetBounds[i].Min)
			}
			appender.AppendChildWidget(item.Widget)
		}
	}
}

func (h *HorizontalLayout) calcItemBounds(context *guigui.Context) {
	h.widgetBounds = slices.Delete(h.widgetBounds, 0, len(h.widgetBounds))

	paddingX, paddingY := formItemPadding(context)

	var x int
	for i, item := range h.items {
		h.widgetBounds = append(h.widgetBounds, image.Rectangle{})

		if item.Widget == nil {
			continue
		}

		var widgetW int
		if item.Widget != nil {
			widgetW, _ = item.Widget.Size(context)
		}
		w := max(widgetW, minFormItemWidth(context))
		baseBounds := guigui.Bounds(h)
		baseBounds.Min.Y += paddingY
		baseBounds.Max.Y -= paddingY
		baseBounds.Min.X += x
		baseBounds.Max.X = baseBounds.Min.X + w

		if item.Widget != nil {
			bounds := baseBounds
			ww, wh := item.Widget.Size(context)
			bounds.Max.Y = bounds.Min.Y + wh
			pX := (w + 2*paddingX - ww) / 2
			if ww < basicwidget.UnitSize(context)+2*paddingX {
				pX = min(pX, max(0, (basicwidget.UnitSize(context)+2*paddingX-ww)/2))
			}
			bounds.Min.X += pX
			bounds.Max.X += pX
			h.widgetBounds[i] = bounds
		}

		x += w + 2*paddingX
	}
}

func (h *HorizontalLayout) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := guigui.Bounds(h)
	bounds.Max.X = bounds.Min.X + h.width(context)
	if !h.background {
		basicwidget.DrawRoundedRect(context, dst, bounds, basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.925), basicwidget.RoundedCornerRadius(context))
	}

	if !h.lineBreak && len(h.items) > 0 {
		paddingX, paddingY := formItemPadding(context)
		x := paddingX
		for _, item := range h.items[:len(h.items)-1] {
			var widgetW int
			if item.Widget != nil {
				widgetW, _ = item.Widget.Size(context)
			}
			w := max(widgetW, minFormItemWidth(context))
			x += w + 2*paddingX

			y0 := float32(bounds.Min.Y + paddingY)
			y1 := float32(bounds.Max.Y - paddingY)
			x := float32(x) + float32(paddingX)
			width := 1 * float32(context.Scale())
			clr := basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.875)
			vector.StrokeLine(dst, x, y0, x, y1, width, clr, false)
		}
	}

	if !h.border {
		basicwidget.DrawRoundedRectBorder(context, dst, bounds, basicwidget.Color(context.ColorMode(), basicwidget.ColorTypeBase, 0.875), basicwidget.RoundedCornerRadius(context), 1*float32(context.Scale()), basicwidget.RoundedRectBorderTypeRegular)
	}
}

func (h *HorizontalLayout) SetHeight(context *guigui.Context, height int) {
	h.heightMinusDefault = height - defaultFormHeight(context)
}

func (h *HorizontalLayout) Size(context *guigui.Context) (int, int) {
	return h.width(context), h.heightMinusDefault + defaultFormHeight(context)
}

func defaultFormHeight(context *guigui.Context) int {
	return 6 * basicwidget.UnitSize(context)
}

func (h *HorizontalLayout) width(context *guigui.Context) int {
	paddingX, _ := formItemPadding(context)

	var x int
	for _, item := range h.items {
		if item.Widget == nil || !guigui.IsVisible(item.Widget) {
			continue
		}
		var widgetW int
		if item.Widget != nil {
			widgetW, _ = item.Widget.Size(context)
		}
		w := max(widgetW, minFormItemWidth(context))
		x += w + 2*paddingX
	}
	return x
}

func minFormItemWidth(context *guigui.Context) int {
	return basicwidget.UnitSize(context)
}
