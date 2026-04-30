// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"image"
	"p86l"
	"p86l/assets"
	"p86l/internal/types"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type About struct {
	guigui.DefaultWidget

	background                 basicwidget.Background
	backButton                 basicwidget.Button
	textPanel                  basicwidget.Panel
	form                       basicwidget.Form
	text1, text2, text3, text4 basicwidget.Text
	image1, image2             aboutIcon

	layoutItems []guigui.LinearLayoutItem
}

func (a *About) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&a.background)
	adder.AddWidget(&a.backButton)
	adder.AddWidget(&a.text1)
	adder.AddWidget(&a.form)
	adder.AddWidget(&a.textPanel)

	v, ok := context.Env(a, modelKeyModel)
	if !ok {
		return nil
	}
	model := v.(*p86l.Model)
	t := model.T()

	context.SetOpacity(&a.background, 0.9)

	a.backButton.SetText("◀")
	a.backButton.OnUp(func(context *guigui.Context) {
		model.SetMode(types.ModeHome)
	})

	a.text1.SetValue(t.Get("about.content"))
	a.text1.SetAutoWrap(true)

	a.text2.SetValue(t.Get("about.lead"))
	a.text2.SetScale(1.2)

	a.text3.SetValue(t.Get("about.dev"))
	a.text3.SetScale(1.2)

	a.text4.SetScale(0.8)
	a.text4.SetAutoWrap(true)
	a.text4.SetMultiline(true)
	a.text4.SetValue(t.Get("about.license"))

	a.image1.setIcon(assets.Images["lead"])
	a.image2.setIcon(assets.Images["dev"])

	a.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &a.text2,
			SecondaryWidget: &a.image1,
		},
		{
			PrimaryWidget:   &a.text3,
			SecondaryWidget: &a.image2,
		},
	})

	a.textPanel.SetContent(&a.text4)
	a.textPanel.SetAutoBorder(true)
	a.textPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	return nil
}

func (a *About) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&a.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	a.layoutItems = slices.Delete(a.layoutItems, 0, len(a.layoutItems))
	a.layoutItems = append(a.layoutItems,
		guigui.LinearLayoutItem{
			Layout: guigui.LinearLayout{
				Direction: guigui.LayoutDirectionHorizontal,
				Items: []guigui.LinearLayoutItem{
					{
						Widget: &a.backButton,
					},
					{
						Size: guigui.FlexibleSize(1),
					},
				},
			},
		},
		guigui.LinearLayoutItem{
			Widget: &a.text1,
		},
		guigui.LinearLayoutItem{
			Widget: &a.form,
		},
		guigui.LinearLayoutItem{
			Widget: &a.textPanel,
			Size:   guigui.FlexibleSize(1),
		},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     a.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

type aboutIcon struct {
	guigui.DefaultWidget

	image basicwidget.Image

	ebitenImage *ebiten.Image
}

func (a *aboutIcon) setIcon(icon *ebiten.Image) {
	a.ebitenImage = icon
}

func (a *aboutIcon) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&a.image)
	a.image.SetImage(a.ebitenImage)

	return nil
}

func (a *aboutIcon) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &a.image,
				Size:   guigui.FixedSize(u * 2),
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (a *aboutIcon) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(u*2, u*2)
}
