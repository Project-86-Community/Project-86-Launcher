// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"p86l"
	"p86l/app"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/fs"
	"p86l/internal/logger"

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

	// Fs
	logger.Info.Println("initialising application filesystem...")
	afs, err := fs.New(*fakeFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: fs init:", err)
		os.Exit(1)
	}

	// Ensure icons are copied to icons directory (skip in fake mode)
	if !*fakeFlag {
		if err := fs.EnsureIcons(afs); err != nil {
			fmt.Fprintln(os.Stderr, "fatal: failed to copy icons:", err)
			os.Exit(1)
		}
	}

	// Logger
	if err := logger.Init(afs, *fakeFlag); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: logger init:", err)
		os.Exit(1)
	}
	logger.LogStartup(VERSION)

	logger.Info.Println("instance lock acquired")

	// Webview
	logger.Info.Println("starting webview thread...")
	wvCh := p86l.RunWebviewThread()
	logger.Info.Println("webview thread ready")

	// Background music
	logger.Info.Println("loading background music...")
	player, err := assets.NewAudioPlayer()
	if err != nil {
		logger.Error.Fatalf("failed to load audio: %v", err)
	}
	defer func() { _ = player.Close() }()

	player.Play()
	logger.Info.Println("background music playing")

	// UI
	logger.Info.Println("building UI root...")
	root := app.NewRoot(afs, wvCh, player)
	if *fakeFlag {
		root.UseFakes(*errorFlag)
		logger.Warn.Println("running in FAKE mode - no real downloads or disk writes")
	}

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
