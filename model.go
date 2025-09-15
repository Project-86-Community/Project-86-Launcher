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
	"bytes"
	"errors"
	"fmt"
	"net"
	"p86l/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Model struct {
	listener net.Listener
	log      LogModel
	file     FileModel
	data     DataModel
	cache    CacheModel

	player *audio.Player
}

// -- Getters for Model --

func (m *Model) Listener() net.Listener {
	return m.listener
}

func (m *Model) Log() *LogModel {
	return &m.log
}

func (m *Model) File() *FileModel {
	return &m.file
}

func (m *Model) Data() *DataModel {
	return &m.data
}

func (m *Model) Cache() *CacheModel {
	return &m.cache
}

func (m *Model) Player() *audio.Player {
	return m.player
}

// -- Setters for Model --

func (m *Model) SetListener(listener net.Listener) {
	m.listener = listener
}

func (m *Model) SetPlayer(player *audio.Player) {
	m.player = player
}

// -- common --

func (m Model) StartBGM() (*audio.Player, error) {
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

func (m *Model) Close() error {
	return errors.Join(m.listener.Close(), m.player.Close(), m.Log().Close(), m.file.Close())
}
