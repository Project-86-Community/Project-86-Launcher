/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2026 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
 * Copyright (C) 2026 Project 86 Community
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

package assets

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/opus"
)

var (
	//go:embed audio/Project_86_OST_legion.opus
	songBytes []byte

	sampleRate   = 48000
	audioContext *audio.Context
)

// NewAudioPlayer decodes the embedded opus file and returns a ready player.
func NewAudioPlayer() (*audio.Player, error) {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	stream, err := opus.DecodeF32(bytes.NewReader(songBytes))
	if err != nil {
		return nil, fmt.Errorf("audio: decode opus: %w", err)
	}

	loop := audio.NewInfiniteLoopF32(stream, stream.Length())

	player, err := audioContext.NewPlayerF32(loop)
	if err != nil {
		return nil, fmt.Errorf("audio: new player: %w", err)
	}

	return player, nil
}
