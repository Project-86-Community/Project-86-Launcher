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
	pd "p86l/internal/debug"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

type ProgressTracker struct {
	model *Model

	filename   string
	totalSize  int64
	downloaded int64
	startTime  time.Time
}

func (p *ProgressTracker) Write(b []byte) (int, error) {
	n := len(b)
	p.downloaded += int64(n)
	p.PrintProgress()
	return n, nil
}

func (p *ProgressTracker) PrintProgress() {
	// Calculate progress.
	currentSize := humanize.Bytes(uint64(p.downloaded))
	totalSize := humanize.Bytes(uint64(p.totalSize))

	// Calculate remaining time.
	elapsed := time.Since(p.startTime).Seconds()
	speed := float64(p.downloaded) / elapsed
	remaining := float64(p.totalSize-p.downloaded) / speed

	remainingTime := humanize.Time(time.Now().Add(time.Duration(remaining) * time.Second))

	// Print the progress.
	output := fmt.Sprintf("\rDownloading %s: %s/%s @ %s/s, %s remaining      ",
		p.filename,
		currentSize,
		totalSize,
		humanize.Bytes(uint64(speed)),
		remainingTime,
	)
	p.model.SetProgress(output)

	// Print newline when done.
	if p.downloaded == p.totalSize {
		log.Info().Str("Downloaded file", p.filename).Msg(pd.NetworkManager)
	}
}
