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
	"io"
	"p86l/internal/debug"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hashicorp/go-getter"
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

func IsNewVersion(currentVersion, newVersion string) (bool, error) {
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

type DownloadType int

const (
	GameDownloadType DownloadType = iota
)

type ProgressTracker struct {
	Model *Model
}

func (p *ProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	return &readProgress{
		Model:      p.Model,
		ReadCloser: stream,
		current:    currentSize,
		total:      totalSize,
		source:     src,
		startTime:  time.Now(),
		lastTime:   time.Now(),
		lastBytes:  currentSize,
	}
}

type readProgress struct {
	Model *Model
	io.ReadCloser
	current   int64
	total     int64
	source    string
	startTime time.Time
	lastTime  time.Time
	lastBytes int64
}

func (r *readProgress) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	if n > 0 {
		now := time.Now()
		r.current += int64(n)

		// Calculate speed and ETA
		timeElapsed := now.Sub(r.lastTime).Seconds()
		if timeElapsed > 0.1 { // Update stats every 100ms to smooth fluctuations
			bytesSinceLast := r.current - r.lastBytes
			speed := float64(bytesSinceLast) / timeElapsed

			remainingBytes := r.total - r.current
			eta := time.Duration(float64(remainingBytes)/speed) * time.Second

			percentage := float64(r.current) / float64(r.total) * 100

			msg := fmt.Sprintf("- Downloading %.0f%% (%s/%s), Speed: %s/s, ETA: %s -",
				percentage,
				humanize.IBytes(uint64(r.current)),
				humanize.IBytes(uint64(r.total)),
				humanize.IBytes(uint64(speed)),
				eta.Round(time.Second),
			)

			r.Model.lunch.SetMsg(msg)

			r.lastTime = now
			r.lastBytes = r.current
		}
	}

	return
}

func DownloadFile(model *Model, dest, sourceUrl string) *debug.Error {
	client := &getter.Client{
		Src:              sourceUrl,
		Dst:              dest,
		Mode:             getter.ClientModeDir,
		ProgressListener: &ProgressTracker{Model: model},
	}

	err := client.Get()
	if err != nil {
		return e.New(err, debug.InternetError, debug.ErrInternetRequestInvalid)
	}

	model.lunch.SetMsg("")
	return nil
}
