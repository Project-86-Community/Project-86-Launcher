// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

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
