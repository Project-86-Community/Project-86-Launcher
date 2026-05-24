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

type Logs struct {
	guigui.DefaultWidget

	background basicwidget.Background
	backButton basicwidget.Button
	titleText  basicwidget.Text
	logPanel   basicwidget.Panel
	logText    basicwidget.Text

	layoutItems []guigui.LinearLayoutItem
}

func (l *Logs) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&l.background)
	adder.AddWidget(&l.backButton)
	adder.AddWidget(&l.titleText)
	adder.AddWidget(&l.logPanel)

	m, ok1 := envMust[*Model](context, l, modelKeyModel)
	ls, ok2 := envMust[service.LogService](context, l, keyLog)
	if !ok1 || !ok2 {
		return nil
	}
	t := m.T()

	context.SetOpacity(&l.background, 0.9)

	l.backButton.SetText("◀")
	l.backButton.OnUp(func(context *guigui.Context) {
		m.SetMode(types.ModeHome)
	})

	l.titleText.SetValue(t.Get("logs.title"))
	l.titleText.SetScale(1.2)

	content, err := ls.ReadLatestLog()
	if err != nil {
		l.logText.SetValue(t.Get("logs.error"))
	} else {
		l.logText.SetValue(content)
	}

	l.logText.SetMultiline(true)
	l.logText.SetWrapMode(basicwidget.WrapModeNormal)

	l.logPanel.SetContent(&l.logText)
	l.logPanel.SetAutoBorder(true)
	l.logPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	return nil
}

func (l *Logs) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&l.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	l.layoutItems = slices.Delete(l.layoutItems, 0, len(l.layoutItems))
	l.layoutItems = append(l.layoutItems,
		guigui.LinearLayoutItem{
			Layout: guigui.LinearLayout{
				Direction: guigui.LayoutDirectionHorizontal,
				Items: []guigui.LinearLayoutItem{
					{Widget: &l.backButton},
					{Size: guigui.FlexibleSize(1)},
				},
			},
		},
		guigui.LinearLayoutItem{Widget: &l.titleText},
		guigui.LinearLayoutItem{Widget: &l.logPanel, Size: guigui.FlexibleSize(1)},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     l.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start: u / 2, Top: u / 2, End: u / 2, Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
