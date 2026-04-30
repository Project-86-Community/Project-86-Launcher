// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"p86l"
	"p86l/internal/types"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type Settings struct {
	guigui.DefaultWidget

	background                basicwidget.Background
	backButton                basicwidget.Button
	form                      basicwidget.Form
	colorModeText             basicwidget.Text
	colorModeSegmentedControl basicwidget.SegmentedControl[string]
	languageText              basicwidget.Text
	languageSelect            basicwidget.Select[language.Tag]
	translateText             basicwidget.Text
	translateToggle           basicwidget.Toggle
	scaleText                 basicwidget.Text
	scaleSegmentedControl     basicwidget.SegmentedControl[float64]
	musicText                 basicwidget.Text
	musicToggle               basicwidget.Toggle

	layoutItems []guigui.LinearLayoutItem
}

func (s *Settings) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&s.background)
	adder.AddWidget(&s.backButton)
	adder.AddWidget(&s.form)

	v, ok := context.Env(s, modelKeyModel)
	if !ok {
		return nil
	}
	model := v.(*p86l.Model)
	t := model.T()

	context.SetOpacity(&s.background, 0.9)

	s.backButton.SetText("◀")
	s.backButton.OnUp(func(context *guigui.Context) {
		model.SetMode(types.ModeHome)
	})

	s.colorModeText.SetValue(t.Get("settings.color_mode"))
	s.colorModeSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[string]{
		{
			Text:  t.Get("settings.color_mode_auto"),
			Value: "",
		},
		{
			Text:  t.Get("settings.color_mode_light"),
			Value: "light",
		},
		{
			Text:  t.Get("settings.color_mode_dark"),
			Value: "dark",
		},
	})
	s.colorModeSegmentedControl.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := s.colorModeSegmentedControl.ItemByIndex(index)
		if !ok {
			context.SetPreferredColorMode(ebiten.ColorModeLight)
			return
		}
		switch item.Value {
		case "light":
			context.SetPreferredColorMode(ebiten.ColorModeLight)
		case "dark":
			context.SetPreferredColorMode(ebiten.ColorModeDark)
		default:
			context.SetPreferredColorMode(ebiten.ColorModeUnknown)
		}
	})
	switch context.PreferredColorMode() {
	case ebiten.ColorModeLight:
		s.colorModeSegmentedControl.SelectItemByValue("light")
	case ebiten.ColorModeDark:
		s.colorModeSegmentedControl.SelectItemByValue("dark")
	default:
		s.colorModeSegmentedControl.SelectItemByValue("")
	}

	s.languageText.SetValue(t.Get("settings.language"))
	s.languageSelect.SetItems([]basicwidget.SelectItem[language.Tag]{
		{
			Text:  "English",
			Value: language.English,
		},
		{
			Text:  "French",
			Value: language.French,
		},
	})
	s.languageSelect.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := s.languageSelect.ItemByIndex(index)
		if !ok {
			context.SetAppLocales(nil)
			model.SetT(language.English.String())
			return
		}
		if item.Value == language.English {
			context.SetAppLocales(nil)
			model.SetT(language.English.String())
			return
		}
		context.SetAppLocales([]language.Tag{item.Value})
		model.SetT(item.Value.String())
	})
	if !s.languageSelect.IsPopupOpen() {
		if locales := context.AppendAppLocales(nil); len(locales) > 0 {
			s.languageSelect.SelectItemByValue(locales[0])
		} else {
			s.languageSelect.SelectItemByValue(language.English)
		}
	}

	// translate changelog or other text that is collected from a API.
	s.translateText.SetValue(t.Get("settings.translate"))
	s.translateToggle.OnValueChanged(func(context *guigui.Context, value bool) {

	})

	s.scaleText.SetValue(t.Get("settings.scale"))
	s.scaleSegmentedControl.SetItems([]basicwidget.SegmentedControlItem[float64]{
		{
			Text:  "50%",
			Value: 0.5,
		},
		{
			Text:  "75%",
			Value: 0.75,
		},
		{
			Text:  "100%",
			Value: 1.0,
		},
		{
			Text:  "125%",
			Value: 1.25,
		},
		{
			Text:  "150%",
			Value: 1.50,
		},
	})
	s.scaleSegmentedControl.OnItemSelected(func(context *guigui.Context, index int) {
		item, ok := s.scaleSegmentedControl.ItemByIndex(index)
		if !ok {
			context.SetAppScale(1)
			return
		}
		context.SetAppScale(item.Value)
	})
	s.scaleSegmentedControl.SelectItemByValue(context.AppScale())

	// disables background music.
	s.musicText.SetValue(t.Get("settings.music"))
	s.musicToggle.OnValueChanged(func(context *guigui.Context, value bool) {
		if value {
			model.PlayBackgroundMusic(true)
		} else {
			model.PlayBackgroundMusic(false)
		}
	})

	s.form.SetItems([]basicwidget.FormItem{
		{
			PrimaryWidget:   &s.colorModeText,
			SecondaryWidget: &s.colorModeSegmentedControl,
		},
		{
			PrimaryWidget:   &s.languageText,
			SecondaryWidget: &s.languageSelect,
		},
		{
			PrimaryWidget:   &s.translateText,
			SecondaryWidget: &s.translateToggle,
		},
		{
			PrimaryWidget:   &s.scaleText,
			SecondaryWidget: &s.scaleSegmentedControl,
		},
		{
			PrimaryWidget:   &s.musicText,
			SecondaryWidget: &s.musicToggle,
		},
	})

	return nil
}

func (s *Settings) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&s.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	s.layoutItems = slices.Delete(s.layoutItems, 0, len(s.layoutItems))
	s.layoutItems = append(s.layoutItems,
		guigui.LinearLayoutItem{
			Layout: guigui.LinearLayout{
				Direction: guigui.LayoutDirectionHorizontal,
				Items: []guigui.LinearLayoutItem{
					{
						Widget: &s.backButton,
					},
					{
						Size: guigui.FlexibleSize(1),
					},
				},
			},
		},
		guigui.LinearLayoutItem{
			Widget: &s.form,
		},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     s.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
