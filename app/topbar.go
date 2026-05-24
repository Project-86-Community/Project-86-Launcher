// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	gocontext "context"
	"fmt"
	"image"
	"p86l/configs"
	"p86l/internal/logger"
	"p86l/internal/service"
	"p86l/internal/types"
	"slices"
	"sync/atomic"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TopBar struct {
	guigui.DefaultWidget

	background     basicwidget.Background
	installButton  basicwidget.Button
	folderMenubar  basicwidget.Menubar[types.Folder]
	settingsButton basicwidget.Button
	helpsMenubar   basicwidget.Menubar[types.Helps]

	dlPopup        basicwidget.Popup
	dlPopupContent guigui.WidgetWithSize[*dlPopupContent]

	downloadDone chan error
	progress     atomic.Pointer[service.InstallProgress]

	layoutItems []guigui.LinearLayoutItem
}

func (t *TopBar) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&t.dlPopup)
	adder.AddWidget(&t.background)
	adder.AddWidget(&t.installButton)
	adder.AddWidget(&t.folderMenubar)
	adder.AddWidget(&t.settingsButton)
	adder.AddWidget(&t.helpsMenubar)

	ds, ok1 := envMust[service.DownloadService](context, t, keyDownload)
	fo, ok2 := envMust[service.FileOpenerService](context, t, keyFileOpen)
	m, ok3 := envMust[*Model](context, t, modelKeyModel)
	wv, ok4 := envMust[chan<- WebviewRequest](context, t, keyWvCh)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil
	}
	_t := m.T()

	t.installButton.SetText(_t.Get("topbar.install"))
	t.installButton.OnUp(func(context *guigui.Context) {
		context.SetEnabled(&t.installButton, false)
		t.downloadDone = make(chan error, 1)
		t.dlPopupContent.Widget().Reset()
		t.dlPopup.SetOpen(true)

		go func() {
			url := configs.Github + "/releases"
			t.dlPopupContent.Widget().SetStatus(_t.Get("topbar.dl.selecting"))
			reply := make(chan string, 1)
			opts := WebviewRequest{Title: _t.Get("topbar.dl.webview"), Source: url, Reply: reply}
			wv <- opts
			url = <-reply

			if url == "" {
				t.downloadDone <- nil
				return
			}

			t.dlPopupContent.Widget().SetStatus(_t.Get("topbar.dl.connecting"))

			err := ds.InstallVersion(
				gocontext.Background(),
				url,
				func(p service.InstallProgress) {
					t.progress.Store(&p)
				},
			)
			t.downloadDone <- err
		}()
	})

	folderItems := []basicwidget.PopupMenuItem[types.Folder]{
		{Text: _t.Get("topbar.f.root"), Value: types.FolderRoot},
		{Text: _t.Get("topbar.f.versions"), Value: types.FolderVersions},
		{Text: _t.Get("topbar.f.logs"), Value: types.FolderLogs},
	}
	t.folderMenubar.SetItems([]basicwidget.MenubarItem{
		{Text: _t.Get("topbar.folders")},
	})
	t.folderMenubar.PopupMenuAt(0).SetItems(folderItems)
	t.folderMenubar.OnItemSelected(func(context *guigui.Context, menuIndex, itemIndex int) {
		if itemIndex < 0 || itemIndex >= len(folderItems) {
			return
		}
		fo.OpenFolder(folderItems[itemIndex].Value)
	})

	t.settingsButton.SetText(_t.Get("topbar.settings"))
	t.settingsButton.OnUp(func(context *guigui.Context) {
		m.SetMode(types.ModeSettings)
	})

	helpsItems := []basicwidget.PopupMenuItem[types.Helps]{
		{Text: _t.Get("topbar.h.cache"), Value: types.HelpsCache},
		{Text: _t.Get("topbar.h.report"), Value: types.HelpsReport},
		{Text: _t.Get("topbar.h.view"), Value: types.HelpsLogs},
		{Border: true},
		{Text: _t.Get("topbar.h.website"), Value: types.HelpsWebsite},
		{Text: "Github", Value: types.HelpsGithub},
		{Text: "Discord", Value: types.HelpsDiscord},
		{Text: "Patreon", Value: types.HelpsPatreon},
		{Border: true},
		{Text: _t.Get("topbar.h.about"), Value: types.HelpsAbout},
	}
	t.helpsMenubar.SetItems([]basicwidget.MenubarItem{
		{Text: _t.Get("topbar.help")},
	})
	t.helpsMenubar.PopupMenuAt(0).SetItems(helpsItems)
	t.helpsMenubar.OnItemSelected(func(context *guigui.Context, menuIndex, itemIndex int) {
		if itemIndex < 0 || itemIndex >= len(helpsItems) {
			return
		}
		switch helpsItems[itemIndex].Value {
		case types.HelpsReport:
			fo.OpenURL(configs.Issues)
		case types.HelpsLogs:
			m.SetMode(types.ModeLogs)
		case types.HelpsWebsite:
			fo.OpenURL(configs.Website)
		case types.HelpsGithub:
			fo.OpenURL(configs.Github)
		case types.HelpsDiscord:
			fo.OpenURL(configs.Discord)
		case types.HelpsPatreon:
			fo.OpenURL(configs.Patreon)
		case types.HelpsAbout:
			m.SetMode(types.ModeAbout)
		}
	})

	t.dlPopupContent.Widget().SetPopup(&t.dlPopup)
	t.dlPopup.SetContent(&t.dlPopupContent)
	t.dlPopup.SetBackgroundDark(true)
	t.dlPopup.SetAnimated(true)

	return nil
}

func (t *TopBar) contentSize(context *guigui.Context) image.Point {
	u := basicwidget.UnitSize(context)
	return image.Pt(int(12*u), int(float64(6.5*float64(u))))
}

func (t *TopBar) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&t.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)
	popupBounds := context.AppBounds()
	t.dlPopup.SetBackgroundBounds(popupBounds)
	contentSize := t.contentSize(context)
	center := image.Point{
		X: popupBounds.Min.X + (popupBounds.Dx()-contentSize.X)/2,
		Y: popupBounds.Min.Y + (popupBounds.Dy()-contentSize.Y)/2,
	}
	layouter.LayoutWidget(&t.dlPopup, image.Rectangle{
		Min: center,
		Max: center.Add(contentSize),
	})

	t.layoutItems = slices.Delete(t.layoutItems, 0, len(t.layoutItems))
	t.layoutItems = append(t.layoutItems,
		guigui.LinearLayoutItem{Widget: &t.installButton},
		guigui.LinearLayoutItem{Widget: &t.folderMenubar},
		guigui.LinearLayoutItem{Widget: &t.settingsButton},
		guigui.LinearLayoutItem{Widget: &t.helpsMenubar},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     t.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start: u / 2, Top: u / 4, End: u / 2, Bottom: u / 4,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (t *TopBar) HandleButtonInput(context *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	m, ok := envMust[*Model](context, t, modelKeyModel)
	if !ok {
		return guigui.HandleInputResult{}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		m.SetMode(types.ModeHome)
		return guigui.HandleInputByWidget(t)
	}

	return guigui.HandleInputResult{}
}

func (t *TopBar) Tick(context *guigui.Context, windowBounds *guigui.WidgetBounds) error {
	m, ok := envMust[*Model](context, t, modelKeyModel)
	if !ok {
		return nil
	}
	_t := m.T()

	if t.downloadDone != nil {
		select {
		case err := <-t.downloadDone:
			t.downloadDone = nil
			t.progress.Store(nil)
			context.SetEnabled(&t.installButton, true)
			if err != nil {
				t.dlPopupContent.Widget().SetPhase(context, types.PhaseError, err.Error())
			} else {
				t.dlPopupContent.Widget().SetPhase(context, types.PhaseDone, _t.Get("topbar.dl.install_complete"))
			}
		default:
			if p := t.progress.Load(); p != nil {
				t.dlPopupContent.Widget().SetInstallProgress(context, *p)
			}
		}
	}
	return nil
}

type dlPopupContent struct {
	guigui.DefaultWidget

	popup *basicwidget.Popup

	titleText   basicwidget.Text
	statusPanel basicwidget.Panel
	statusText  basicwidget.Text
	progress    basicwidget.Slider
	closeButton basicwidget.Button

	showClose   atomic.Bool
	layoutItems []guigui.LinearLayoutItem
}

func (d *dlPopupContent) Reset() {
	d.titleText.SetValue("")
	d.progress.SetValue(0)
	d.statusText.SetValue("")
	d.showClose.Store(false)
}

func (d *dlPopupContent) SetPopup(popup *basicwidget.Popup) {
	d.popup = popup
}

func (d *dlPopupContent) SetInstallProgress(context *guigui.Context, p service.InstallProgress) {
	m, ok := envMust[*Model](context, d, modelKeyModel)
	if !ok {
		return
	}
	t := m.T()

	switch p.Phase {
	case types.PhaseDownload:
		d.titleText.SetValue(t.Get("topbar.dl.downloading"))
		if p.DownloadTotal > 0 {
			pct := float64(p.DownloadDone) / float64(p.DownloadTotal)
			d.progress.SetValue(int(pct * 100))
			d.statusText.SetValue(fmt.Sprintf("%.0f%%  (%s / %s)",
				pct*100,
				formatBytes(p.DownloadDone),
				formatBytes(p.DownloadTotal),
			))
		} else {
			d.statusText.SetValue(t.Get("topbar.dl.connecting"))
		}
	case types.PhaseExtract:
		d.titleText.SetValue(t.Get("topbar.dl.extracting"))
		if p.ExtractTotal > 0 {
			pct := float64(p.ExtractDone) / float64(p.ExtractTotal)
			d.progress.SetValue(int(pct * 100))
			d.statusText.SetValue(fmt.Sprintf(t.Get("topbar.dl.extracting_files"), p.ExtractDone, p.ExtractTotal))
		} else {
			d.progress.SetValue(0)
			d.statusText.SetValue(fmt.Sprintf(t.Get("topbar.dl.extracted_count"), p.ExtractDone))
		}
	}
}

func (d *dlPopupContent) SetPhase(context *guigui.Context, phase types.Phase, msg string) {
	m, ok := envMust[*Model](context, d, modelKeyModel)
	if !ok {
		return
	}
	t := m.T()

	switch phase {
	case types.PhaseDone:
		logger.Info.Printf("install done: %s", msg)
		d.titleText.SetValue(t.Get("topbar.dl.done"))
		d.progress.SetValue(100)
		d.statusText.SetValue(msg)
		d.showClose.Store(true)
	case types.PhaseError:
		logger.Error.Printf("install failed: %v", msg)
		d.titleText.SetValue(t.Get("topbar.dl.failed"))
		d.progress.SetValue(0)
		d.statusText.SetValue(msg)
		d.showClose.Store(true)
	}
}

func (d *dlPopupContent) SetStatus(msg string) {
	d.statusText.SetValue(msg)
}

func (d *dlPopupContent) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&d.titleText)
	adder.AddWidget(&d.progress)
	adder.AddWidget(&d.statusPanel)
	adder.AddWidget(&d.closeButton)

	m, ok := envMust[*Model](context, d, modelKeyModel)
	if !ok {
		return nil
	}
	t := m.T()

	d.titleText.SetScale(1.2)

	d.progress.SetMinimumValue(0)
	d.progress.SetMaximumValue(100)
	context.SetEnabled(&d.progress, false)

	d.statusText.SetWrapMode(basicwidget.WrapModeNormal)

	d.statusPanel.SetContent(&d.statusText)
	d.statusPanel.SetAutoBorder(true)
	d.statusPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	showClose := d.showClose.Load()
	d.closeButton.SetText(t.Get("common.close"))
	context.SetEnabled(&d.closeButton, showClose)

	d.closeButton.OnUp(func(context *guigui.Context) {
		if d.popup != nil {
			d.popup.SetOpen(false)
			d.showClose.Store(false)
		}
	})

	return nil
}

func (d *dlPopupContent) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	d.layoutItems = slices.Delete(d.layoutItems, 0, len(d.layoutItems))
	d.layoutItems = append(d.layoutItems,
		guigui.LinearLayoutItem{Widget: &d.titleText},
		guigui.LinearLayoutItem{Widget: &d.progress},
		guigui.LinearLayoutItem{Widget: &d.statusPanel, Size: guigui.FixedSize(int(float64(u) * 1.5))},
		guigui.LinearLayoutItem{Widget: &d.closeButton},
	)
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     d.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start: u / 2, Top: u / 4, End: u / 2, Bottom: u / 4,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
