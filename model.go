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

	mode string

	rateLimit RateLimitModel
	play      PlayModel
	app       AppModel
	data      DataModel
	cache     CacheModel
}

// -- new --
func (m *Model) Listener() net.Listener {
	return m.listener
}

func (m *Model) Log() *LogModel {
	return &m.log
}

func (m *Model) SetListener(listener net.Listener) {
	m.listener = listener
}

// -- new - common --

func (m *Model) Close() error {
	return errors.Join(m.listener.Close(), m.Log().Close())
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

func (m *Model) GameExecutablePath() string {
	am := m.App()
	data := m.Data()
	if data.File().UsePreRelease {
		return am.FileSystem().PathFilePrerelease()
	}
	return am.FileSystem().PathFileGame()
}

// -- Models --

func (m *Model) App() *AppModel {
	return &m.app
}

func (m *Model) RateLimits() *RateLimitModel {
	return &m.rateLimit
}

func (m *Model) Play() *PlayModel {
	return &m.play
}

func (m *Model) Data() *DataModel {
	return &m.data
}

func (m *Model) Cache() *CacheModel {
	return &m.cache
}

// -- RatelimitModel --

type RateLimitModel struct {
	progress bool
}

func (r *RateLimitModel) Progress() bool {
	return r.progress
}

func (r *RateLimitModel) SetProgress(value bool) {
	r.progress = value
}
