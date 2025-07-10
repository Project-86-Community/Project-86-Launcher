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
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
)

type DownloadProgressTracker struct {
	Model      *Model
	total      int64
	current    int64
	startTime  time.Time
	filename   string
	lastUpdate time.Time
}

func (pt *DownloadProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
	pt.total = totalSize
	pt.startTime = time.Now()
	pt.lastUpdate = time.Now()

	// Extract filename from source URL.
	pt.filename = filepath.Base(src)

	return &progressReader{
		ReadCloser: stream,
		tracker:    pt,
	}
}

type progressReader struct {
	io.ReadCloser
	tracker *DownloadProgressTracker
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.ReadCloser.Read(p)
	pr.tracker.current += int64(n)

	// Update progress every 100ms.
	if time.Since(pr.tracker.lastUpdate) >= 100*time.Millisecond {
		pr.tracker.printProgress()
		pr.tracker.lastUpdate = time.Now()
	}

	if err == io.EOF {
		pr.tracker.printProgress()
	}

	return n, err
}

func (pt *DownloadProgressTracker) printProgress() {
	if pt.total == 0 {
		return
	}

	// Calculate download speed and remaining time.
	elapsed := time.Since(pt.startTime)
	if elapsed == 0 {
		elapsed = 1 * time.Millisecond
	}

	bytesPerSecond := float64(pt.current) / elapsed.Seconds()
	remainingBytes := pt.total - pt.current

	// Calculate estimated completion time.
	var remainingTime string
	if bytesPerSecond > 0 && remainingBytes > 0 {
		remainingSeconds := float64(remainingBytes) / bytesPerSecond
		estimatedCompletion := time.Now().Add(time.Duration(remainingSeconds * float64(time.Second)))

		remainingTime = humanize.Time(estimatedCompletion)

		// Remove "in " prefix if present for cleaner output.
		if len(remainingTime) > 3 && remainingTime[:3] == "in " {
			remainingTime = remainingTime[3:] + " remaining"
		} else {
			remainingTime = remainingTime + " remaining"
		}
	} else {
		remainingTime = "calculating..."
	}

	currentSize := humanize.Bytes(uint64(pt.current))
	totalSize := humanize.Bytes(uint64(pt.total))
	percentage := float64(pt.current) / float64(pt.total) * 100

	output := fmt.Sprintf("\rDownloading %s: %s/%s (%.1f%%), %s",
		pt.filename, currentSize, totalSize, percentage, remainingTime)
	pt.Model.SetProgress(output)
}

