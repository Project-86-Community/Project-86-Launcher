// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/fyne-io/image/ico"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images/*
var imagesFS embed.FS

var (
	// Images holds all jpg/png files as *ebiten.Image, keyed by filename without extension.
	// e.g. "images/logo.png" → Images["logo"]
	Images map[string]*ebiten.Image

	// Icons holds all ico files as []image.Image.
	// e.g. "images/app.ico" → Icons["app"]
	Icons map[string][]image.Image

	// IconData holds raw bytes of icon files, keyed by filename with extension.
	// e.g. "images/icon.ico" → IconData["icon.ico"]
	IconData map[string][]byte
)

func init() {
	Images = make(map[string]*ebiten.Image)
	Icons = make(map[string][]image.Image)
	IconData = make(map[string][]byte)

	entries, err := imagesFS.ReadDir("images")
	if err != nil {
		panic(fmt.Sprintf("assets: failed to read images dir: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		key := strings.TrimSuffix(name, filepath.Ext(name))
		path := "images/" + name

		data, err := imagesFS.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("assets: failed to read %s: %v", path, err))
		}

		switch ext {
		case ".jpg", ".jpeg", ".png":
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(fmt.Sprintf("assets: failed to decode image %s: %v", path, err))
			}
			Images[key] = ebiten.NewImageFromImage(img)

		case ".ico":
			imgs, err := ico.DecodeAll(bytes.NewReader(data))
			if err != nil {
				panic(fmt.Sprintf("assets: failed to decode ico %s: %v", path, err))
			}
			Icons[key] = imgs
		}

		if strings.HasPrefix(name, "icon.") {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".ico" || ext == ".icns" || ext == ".png" {
				IconData[name] = data
			}
		}
	}
}
