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
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

type ProgressTracker struct {
	model       *Model
	filename    string
	totalSize   int64
	currentSize int64
	startTime   time.Time
	lastTime    time.Time
	lastSize    int64
	lastPrint   time.Time
	mu          sync.Mutex
}

func (p *ProgressTracker) Write(data []byte) (int, error) {
	p.mu.Lock()
	p.currentSize += int64(len(data))

	if time.Since(p.lastPrint) >= 500*time.Millisecond {
		p.PrintProgress()
		p.lastPrint = time.Now()
	}

	p.mu.Unlock()
	return len(data), nil
}

func (p *ProgressTracker) PrintProgress() {
	// Calculate progress.
	currentSize := humanize.Bytes(uint64(p.currentSize))
	totalSize := humanize.Bytes(uint64(p.totalSize))

	now := time.Now()
	var speed float64

	if !p.lastTime.IsZero() {
		elapsed := now.Sub(p.lastTime).Seconds()
		if elapsed > 0 {
			speed = float64(p.currentSize-p.lastSize) / elapsed
		}
	}

	p.lastTime = now
	p.lastSize = p.currentSize

	var remainingStr string
	if speed > 0 {
		remaining := float64(p.totalSize-p.currentSize) / speed
		remainingDuration := time.Duration(remaining) * time.Second
		remainingStr = humanize.RelTime(time.Now(), time.Now().Add(remainingDuration), "remaining", "ago")
	} else {
		remainingStr = "calculating..."
	}

	// Print the progress.
	output := fmt.Sprintf("Downloading %s: %s/%s @ %s/s, %s",
		p.filename,
		currentSize,
		totalSize,
		humanize.Bytes(uint64(speed)),
		remainingStr,
	)
	p.model.SetProgress(output)

	// Print newline when done.
	if p.currentSize == p.totalSize {
		log.Info().Str("Downloaded file", p.filename).Msg(pd.NetworkManager)
	}
}
