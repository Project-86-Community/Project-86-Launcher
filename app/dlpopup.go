// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"fmt"
	"p86l/internal/logger"
	"p86l/internal/service"
	"p86l/internal/types"
	"slices"
	"sync/atomic"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

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
