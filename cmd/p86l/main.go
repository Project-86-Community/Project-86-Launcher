// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"p86l/app"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/service"

	"github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"
)

const VERSION = "dev"

func main() {
	fakeFlag := flag.Bool("fake", false, "use fake downloader/extractor for UI testing")
	errorFlag := flag.Bool("error", false, "simulate errors in fake mode (requires --fake)")
	flag.Parse()

	// Prevent multiple instances.
	logger.Info.Println("acquiring instance lock...")
	lock, err := tryLock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: another instance is already running")
		os.Exit(1)
	}
	defer func() { _ = lock.Close() }()

	logger.Info.Println("initialising application filesystem...")
	afs, err := fs.New(*fakeFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: fs init:", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Close() }()

	// Ensure icons are copied to icons directory (skip in fake mode)
	if !*fakeFlag {
		if err := fs.EnsureIcons(afs); err != nil {
			fmt.Fprintln(os.Stderr, "fatal: failed to copy icons:", err)
			os.Exit(1)
		}
	}

	if err := logger.Init(afs, *fakeFlag); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: logger init:", err)
		os.Exit(1)
	}
	logger.LogStartup(VERSION)

	logger.Info.Println("instance lock acquired")

	logger.Info.Println("starting webview thread...")
	var wvCh chan<- app.WebviewRequest
	if *fakeFlag {
		wvCh = app.FakeWebviewThread()
		logger.Info.Println("webview thread ready (fake)")
	} else {
		wvCh = app.RunWebviewThread()
		logger.Info.Println("webview thread ready")
	}

	logger.Info.Println("loading background music...")
	player, err := assets.NewAudioPlayer()
	if err != nil {
		logger.Error.Fatalf("failed to load audio: %v", err)
	}
	defer func() { _ = player.Close() }()

	player.Play()
	logger.Info.Println("background music playing")

	var (
		download service.DownloadService
		launch   service.LaunchService
		fileOpen service.FileOpenerService
		logSvc   service.LogService
		shortcut service.ShortcutService
		version  service.VersionService
	)

	if *fakeFlag {
		fakeDS := service.NewFakeDownloadService(afs)
		fakeDS.FakeError = *errorFlag
		download = fakeDS
		launch = service.NewFakeLaunchService()
		fileOpen = service.NewFakeFileOpenerService()
		logSvc = service.NewFakeLogService()
		shortcut = service.NewFakeShortcutService()
		version = service.NewFakeVersionService()
		logger.Warn.Println("running in FAKE mode")
	} else {
		download = service.NewDownloadService(afs, fs.GrabDownloader{}, fs.FastExtractor{})
		launch = service.NewLaunchService()
		fileOpen = service.NewFileOpenerService(afs)
		logSvc = service.NewLogService(afs)
		shortcut = service.NewShortcutService(afs)
		version = service.NewVersionService(afs)
	}

	logger.Info.Println("building UI root...")
	model := app.NewModel()
	root := app.NewRoot(download, launch, fileOpen, logSvc, shortcut, version, model, player, wvCh)

	logger.Info.Printf("setting window title: %s %s", configs.Title, VERSION)
	ebiten.SetWindowIcon(assets.Icons["icon"])

	op := &guigui.RunOptions{
		Title:         fmt.Sprintf("%s %s", configs.Title, VERSION),
		WindowMinSize: configs.WindowMinSize,
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}

	logger.Info.Println("entering main UI loop")
	if err := guigui.Run(root, op); err != nil {
		logger.Error.Fatalf("UI loop exited with error: %v", err)
	}
	logger.Info.Println("UI loop exited cleanly - shutting down")
}
