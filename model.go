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
	"p86l/internal/log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/solarlune/resound"
)

type Model struct {
	listener net.Listener
	log      LogModel
	file     FileModel
	data     DataModel
	cache    CacheModel
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

// -- Setters for Model --

func (m *Model) SetListener(listener net.Listener) {
	m.listener = listener
}

// -- common --

func (m Model) StartBGM() {
	const sampleRate = 44100
	audio.NewContext(sampleRate)
	reader := bytes.NewReader(assets.P86lOst)

	stream, err := vorbis.DecodeWithSampleRate(sampleRate, reader)
	if err != nil {
		m.Log().Logger().Err(fmt.Errorf("vorbis decoder: %w", err)).Msg(log.ErrorManager.String())
		return
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())

	player, err := resound.NewPlayer("bgm", loop)
	if err != nil {
		m.Log().Logger().Err(fmt.Errorf("new player: %w", err)).Msg(log.ErrorManager.String())
		return
	}

	player.Play()
}

func (m *Model) Close() error {
	return errors.Join(m.listener.Close(), m.Log().Close(), m.file.Close())
}
