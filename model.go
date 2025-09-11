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
	"errors"
	"net"
	"p86l/internal/log"
)

type Model struct {
	listener net.Listener
	log      LogModel
	file     FileModel
	update   *UpdateModel

	mode string
}

// -- new --
func (m *Model) Listener() net.Listener {
	return m.listener
}

func (m *Model) Log() *LogModel {
	return &m.log
}

func (m *Model) File() *FileModel {
	return &m.file
}

func (m *Model) Update() *UpdateModel {
	if m.update == nil {
		m.update = &UpdateModel{
			urlChan: make(chan string, 10),
		}
	}
	return m.update
}

func (m *Model) SetListener(listener net.Listener) {
	m.listener = listener
}

// -- new - common --

func (m *Model) Close() error {
	return errors.Join(m.listener.Close(), m.Log().Close(), m.file.Close())
}

// -- Getters for Model --

func (m *Model) Mode() string {
	if m.mode == "" {
		return "home"
	}
	return m.mode
}

// -- Setters for Model --

func (m *Model) SetMode(mode string) {
	m.Log().logger.Info().Str("Page", mode).Msg(log.AppManager.String())
	m.mode = mode
}
