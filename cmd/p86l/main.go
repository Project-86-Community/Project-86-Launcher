package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"p86l"
	"p86l/app"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/fs"

	"github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fake := flag.Bool("fake", false, "use fake downloader/extractor for UI testing")
	flag.Parse()

	// Prevent multiple instances.
	lock, err := tryLock()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = lock.Close() }()

	// Fs
	fs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Webview
	wvCh := p86l.RunWebviewThread()

	// Background music
	player, err := assets.NewAudioPlayer()
	if err != nil {
		log.Fatal(err)
	}
	player.Play()

	root := app.NewRoot(fs, wvCh, player)
	if *fake {
		root.UseFakes()
	}

	ebiten.SetWindowIcon(assets.Icons["icon"])
	op := &guigui.RunOptions{
		Title:         configs.Title,
		WindowMinSize: configs.WindowMinSize,
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(root, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
