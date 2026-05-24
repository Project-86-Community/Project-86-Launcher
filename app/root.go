// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"image"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/logger"
	"p86l/internal/service"
	"p86l/internal/types"
	"slices"
	"sync"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/text/language"
)

var (
	keyDownload   = guigui.GenerateEnvKey()
	keyLaunch     = guigui.GenerateEnvKey()
	keyFileOpen   = guigui.GenerateEnvKey()
	keyLog        = guigui.GenerateEnvKey()
	keyShortcut   = guigui.GenerateEnvKey()
	keyVersion    = guigui.GenerateEnvKey()
	keyPlayer     = guigui.GenerateEnvKey()
	keyWvCh       = guigui.GenerateEnvKey()
	modelKeyModel = guigui.GenerateEnvKey()
)

type Root struct {
	guigui.DefaultWidget

	backgroundImage basicwidget.Image
	topbar          TopBar
	home            Home
	settings        Settings
	about           About
	logs            Logs
	bottombar       BottomBar

	download service.DownloadService
	launch   service.LaunchService
	fileOpen service.FileOpenerService
	logSvc   service.LogService
	shortcut service.ShortcutService
	version  service.VersionService
	model    *Model
	player   *audio.Player
	wvCh     chan<- WebviewRequest

	initOnce sync.Once

	layoutItems []guigui.LinearLayoutItem
}

func NewRoot(
	download service.DownloadService,
	launch service.LaunchService,
	fileOpen service.FileOpenerService,
	logSvc service.LogService,
	shortcut service.ShortcutService,
	version service.VersionService,
	model *Model,
	player *audio.Player,
	wvCh chan<- WebviewRequest,
) *Root {
	return &Root{
		download: download,
		launch:   launch,
		fileOpen: fileOpen,
		logSvc:   logSvc,
		shortcut: shortcut,
		version:  version,
		model:    model,
		player:   player,
		wvCh:     wvCh,
	}
}

func (r *Root) Env(context *guigui.Context, key guigui.EnvKey, source *guigui.EnvSource) (any, bool) {
	switch key {
	case keyDownload:
		return r.download, true
	case keyLaunch:
		return r.launch, true
	case keyFileOpen:
		return r.fileOpen, true
	case keyLog:
		return r.logSvc, true
	case keyShortcut:
		return r.shortcut, true
	case keyVersion:
		return r.version, true
	case modelKeyModel:
		return r.model, true
	case keyPlayer:
		return r.player, true
	case keyWvCh:
		return r.wvCh, true
	default:
		return nil, false
	}
}

func (r *Root) contentWidget() guigui.Widget {
	switch r.model.Mode() {
	case types.ModeHome:
		return &r.home
	case types.ModeSettings:
		return &r.settings
	case types.ModeAbout:
		return &r.about
	case types.ModeLogs:
		return &r.logs
	}
	return nil
}

func (r *Root) handleBackgroundImage(widgetBounds *guigui.WidgetBounds) image.Rectangle {
	imgWidth := assets.Images["banner"].Bounds().Dx()
	imgHeight := assets.Images["banner"].Bounds().Dy()

	windowBounds := widgetBounds.Bounds()
	windowWidth := windowBounds.Dx()
	windowHeight := windowBounds.Dy()

	imgAspectRatio := float64(imgWidth) / float64(imgHeight)
	windowAspectRatio := float64(windowWidth) / float64(windowHeight)

	var newWidth, newHeight int
	var xOffset, yOffset int

	if imgAspectRatio > windowAspectRatio {
		// The image is wider than the window. Scale by height and crop width.
		newHeight = windowHeight
		newWidth = int(float64(windowHeight) * imgAspectRatio)
		xOffset = 0
		yOffset = 0
	} else {
		// The image is taller than the window. Scale by width and crop height.
		newWidth = windowWidth
		newHeight = int(float64(windowWidth) / imgAspectRatio)
		xOffset = 0
		yOffset = (windowHeight - newHeight) / 2
	}

	backgroundImagePosition := image.Pt(xOffset, yOffset)
	backgroundImageSize := image.Pt(newWidth, newHeight)

	return image.Rect(
		backgroundImagePosition.X,
		backgroundImagePosition.Y,
		backgroundImagePosition.X+backgroundImageSize.X,
		backgroundImagePosition.Y+backgroundImageSize.Y,
	)
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	r.initOnce.Do(func() {
		var ebdi ebiten.DebugInfo
		ebiten.ReadDebugInfo(&ebdi)
		logger.Info.Printf("Graphics library %s", ebdi.GraphicsLibrary.String())

		locale := context.FirstLocale()
		logger.Info.Printf("First locale: %s", locale.String())
		if locale == (language.Tag{}) {
			t := assets.NewT(language.English.String())
			r.model.SetT(&t)
		} else {
			t := assets.NewT(locale.String())
			r.model.SetT(&t)
		}
	})

	adder.AddWidget(&r.backgroundImage)
	adder.AddWidget(&r.topbar)
	if content := r.contentWidget(); content != nil {
		adder.AddWidget(content)
	}
	adder.AddWidget(&r.bottombar)

	r.backgroundImage.SetImage(assets.Images["banner"])

	{
		w := int(float64(configs.WindowMinSize.X)*context.AppScale()) + basicwidget.UnitSize(context)*2
		h := int(float64(configs.WindowMinSize.Y)*context.AppScale()) + basicwidget.UnitSize(context)*2
		context.SetWindowSizeLimits(w, h, -1, -1)
	}

	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &r.backgroundImage,
				Size:   guigui.FlexibleSize(1),
			},
		},
	}).LayoutWidgets(context, r.handleBackgroundImage(widgetBounds), layouter)

	u := basicwidget.UnitSize(context)
	r.layoutItems = slices.Delete(r.layoutItems, 0, len(r.layoutItems))
	if r.model.Mode() == types.ModeHome {
		r.layoutItems = append(r.layoutItems,
			guigui.LinearLayoutItem{
				Widget: &r.topbar,
				Size:   guigui.FixedSize(int(float64(u) * 1.6)),
			},
			guigui.LinearLayoutItem{
				Widget: r.contentWidget(),
				Size:   guigui.FlexibleSize(1),
			},
			guigui.LinearLayoutItem{
				Widget: &r.bottombar,
				Size:   guigui.FixedSize(u),
			},
		)
	} else {
		r.layoutItems = append(r.layoutItems,
			guigui.LinearLayoutItem{
				Widget: r.contentWidget(),
				Size:   guigui.FlexibleSize(1),
			},
		)
	}
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     r.layoutItems,
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}
