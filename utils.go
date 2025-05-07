/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
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
	"fmt"
	"p86l/internal/debug"
	"strings"

	version "github.com/hashicorp/go-version"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
)

func SetLanguage(lang string) {
	lLocalizer = i18n.NewLocalizer(lBundle, lang)
}

func T(key string) string {
	lMsg, err := lLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return fmt.Sprintf("!{%s}", key)
	}

	return lMsg
}

func OpenBrowser(url string) {
	log.Info().Str("Url", url).Msg("OpenBrowser")
	if err := browser.OpenURL(url); err != nil {
		e.SetPopup(e.New(err, debug.InternetError, debug.ErrBrowserOpen))
	}
}

func IsValidGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		!strings.Contains(filename, "dev")
}

func CheckNewerVersion(currentVersion, newVersion string) (bool, error) {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, fmt.Errorf("invalid current version: %w", err)
	}

	newer, err := version.NewVersion(newVersion)
	if err != nil {
		return false, fmt.Errorf("invalid new version: %w", err)
	}

	return newer.GreaterThan(current), nil
}

// func DownloadFile(model *Model, url, filepath string) *debug.Error {
// 	client := grab.NewClient()
// 	req, err := grab.NewRequest(filepath, url)
// 	if err != nil {
// 		return e.New(err, debug.InternetError, debug.ErrNewRequest)
// 	}
//
// 	resp := client.Do(req)
//
// 	t := time.NewTicker(500 * time.Millisecond)
// 	defer t.Stop()
//
// Loop:
// 	for {
// 		select {
// 		case <-t.C:
// 			speed := float64(resp.BytesPerSecond()) / 1024 / 1024 // Speed in MB/s
// 			eta := resp.ETA()
// 			etaStr := human_duration.String(eta.Sub(time.Now()), "second")
//
// 			model.play.SetDownloadMsg(fmt.Sprintf("(%.2f%%) %.2f, ETA: %s", 100*resp.Progress(), speed, etaStr))
// 		case <-resp.Done:
// 			break Loop
// 		}
// 	}
//
// 	if err := resp.Err(); err != nil {
// 		return e.New(err, debug.InternetError, debug.ErrFailedDownload)
// 	}
//
// 	return nil
// }
