/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
 * Copyright (C) 2025 Project 86 Community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package p86l

import (
	"bytes"
	"cmp"
	"fmt"
	"image"
	"os"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/github"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fyne-io/image/ico"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hashicorp/go-version"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	PathGameStable     = filepath.Join(configs.FolderBuild, configs.FolderGame, configs.FileGame)
	PathGamePreRelease = filepath.Join(configs.FolderBuild, configs.FolderPreRelease, configs.FileGame)
)

func GetIcons() ([]image.Image, error) {
	images, err := ico.DecodeAll(bytes.NewReader(assets.P86lIco))
	if err != nil {
		return nil, fmt.Errorf("failed to decode icons: %w", err)
	}

	return images, nil
}

func GetUsername() string {
	var username string
	switch runtime.GOOS {
	case "windows":
		username = os.Getenv("USERNAME")
	default:
		username = os.Getenv("USER")
	}

	if username == "" {
		username = os.Getenv("LOGNAME")
	}
	return strings.TrimSpace(username)
}

func MergeRectangles(rects ...image.Rectangle) image.Rectangle {
	if len(rects) == 0 {
		return image.Rectangle{}
	}

	mergedMin := rects[0].Min
	mergedMax := rects[0].Max

	for _, r := range rects[1:] {
		if r.Min.X < mergedMin.X {
			mergedMin.X = r.Min.X
		}
		if r.Min.Y < mergedMin.Y {
			mergedMin.Y = r.Min.Y
		}
		if r.Max.X > mergedMax.X {
			mergedMax.X = r.Max.X
		}
		if r.Max.Y > mergedMax.Y {
			mergedMax.Y = r.Max.Y
		}
	}

	return image.Rectangle{
		Min: mergedMin,
		Max: mergedMax,
	}
}

func NewBGMPlayer() (*audio.Player, error) {
	const sampleRate = 44100
	actx := audio.NewContext(sampleRate)

	reader := bytes.NewReader(assets.P86lOst)
	stream, err := vorbis.DecodeF32(reader)
	if err != nil {
		return nil, fmt.Errorf("vorbis decoder: %w", err)
	}

	loop := audio.NewInfiniteLoopF32(stream, stream.Length())

	player, err := actx.NewPlayerF32(loop)
	if err != nil {
		return nil, fmt.Errorf("new player: %w", err)
	}

	return player, nil
}

func T(key string) string {
	keyMsg, err := assets.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return fmt.Sprintf("!%s", key) // Fallback to key if translation fails.
	}
	return keyMsg
}

func FormattedCacheExpireText(cache CacheFile) string {
	if cache.RateLimit != nil {
		var endPart string
		if cache.RateLimit.Remaining < 60 {
			endPart = humanize.RelTime(time.Now(), time.Unix(cache.RateLimit.Reset, 0), "remaining", "ago")
		}

		return fmt.Sprintf(
			"%d / %d - requests \n %s",
			cache.RateLimit.Remaining,
			cache.RateLimit.Limit,
			endPart,
		)
	}

	return "..."
}

func ReleasesChangelogText(cache CacheFile, usePreRelease bool) string {
	if cache.Releases != nil {
		if usePreRelease {
			return fmt.Sprintf("%s\n\n%s", cache.Releases.PreRelease.Name, cache.Releases.PreRelease.Body)
		}
		return fmt.Sprintf("%s\n\n%s", cache.Releases.Stable.Name, cache.Releases.Stable.Body)
	}

	return "..."
}

func GameVersionText(cache CacheFile, usePreRelease bool) string {
	if cache.Releases != nil {
		if usePreRelease {
			return cache.Releases.PreRelease.TagName
		}
		return cache.Releases.Stable.TagName
	}

	return "..."
}

func ReleasesDownloadCountText(cache CacheFile, usePreRelease bool) string {
	if cache.Releases != nil {
		if usePreRelease {
			_, debugGameAsset := GetAssets(cache.Releases.PreRelease.Assets)
			return fmt.Sprintf("%d", debugGameAsset.DownloadCount)
		}
		gameAsset, _ := GetAssets(cache.Releases.Stable.Assets)
		return fmt.Sprintf("%d", gameAsset.DownloadCount)
	}

	return "..."
}

func isGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		!strings.Contains(filename, "dev")
}

func isDebugGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		strings.Contains(filename, "dev")
}

func GetAssets(assets []github.ReleaseAsset) (*github.ReleaseAsset, *github.ReleaseAsset) {
	var gameAsset *github.ReleaseAsset
	var debugGameAsset *github.ReleaseAsset

	for _, asset := range assets {
		switch {
		case isGameFile(asset.Name):
			gameAsset = &asset
			continue
		case isDebugGameFile(asset.Name):
			debugGameAsset = &asset
			continue
		}
	}

	return gameAsset, debugGameAsset
}

// current, new
func IsNewVersion(v1, v2 string) (bool, error) {
	c, err1 := version.NewVersion(v1)
	n, err2 := version.NewVersion(v2)

	if err := cmp.Or(err1, err2); err != nil {
		return false, err
	}

	return n.GreaterThan(c), nil
}
