/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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
	"encoding/json"
	"p86l/internal/log"
	"sync"

	"golang.org/x/text/language"
)

type SidebarPage int

const (
	PageHome SidebarPage = iota
	PagePlay
	PageSettings
	PageAbout
)

type DataRemember struct {
	WSizeX int  `json:"wsizex"`
	WSizeY int  `json:"wsizey"`
	WPosX  int  `json:"wposx"`
	WPosY  int  `json:"wposy"`
	Page   int  `json:"page"`
	Active bool `json:"active"`
}

type DataFile struct {
	Lang               string       `json:"lang"`
	TranslateChangelog bool         `json:"translate_changelog"`
	UseDarkmode        bool         `json:"use_darkmode"`
	AppScale           float64      `json:"app_scale"`
	DisableBgMusic     bool         `json:"disable_bgm"`
	UsePreRelease      bool         `json:"use_pre_release"`
	Remember           DataRemember `json:"remember"`
}

type Data struct {
	mu   sync.RWMutex
	file DataFile
}

func NewData(initial DataFile) *Data {
	return &Data{
		file: initial,
	}
}

func (d *Data) Lang() (language.Tag, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tag, err := language.Parse(d.file.Lang)
	if err != nil {
		return language.English, err
	}

	return tag, nil
}

func (d *Data) Get() DataFile {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.file
}

func (d *Data) Update(fn func(*DataFile)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	fn(&d.file)
}

func (m *Model) loadData() error {
	if !m.fs.Exist(m.dataPath) {
		m.logger.Info().Str(log.Lifecycle, "data file does not exist, using defaults").Msg(log.FileManager.String())

		m.data.Update(func(df *DataFile) {
			df.Lang = language.English.String()
			df.UseDarkmode = false
			df.AppScale = 1
			df.DisableBgMusic = false
			df.UsePreRelease = false
		})
		m.isNew = true

		return nil
	}

	jsonData, err := m.fs.Load(m.dataPath)
	if err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to load data").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	var data DataFile
	if err := json.Unmarshal(jsonData, &data); err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to unmarshal data").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	m.data.Update(func(df *DataFile) {
		*df = data
	})

	m.logger.Info().Str(log.Lifecycle, "data loaded successfully").Any("data", data).Msg(log.FileManager.String())
	return nil
}

func (m *Model) saveData() error {
	data := m.data.Get()

	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to marshal data").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	if err := m.fs.Save(m.dataPath, jsonData); err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to save data").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	m.logger.Info().Str(log.Lifecycle, "data saved successfully").Msg(log.FileManager.String())
	return nil
}

func (m *Model) Data() *Data {
	return m.data
}

func (m *Model) SyncData() error {
	if m.syncDataFn != nil {
		return m.syncDataFn(m, m.isNew)
	}
	return nil
}

// -- commands --

type ResetDataCommand struct{}

func (r ResetDataCommand) Execute(m *Model) {
	m.data.Update(func(df *DataFile) {
		df.Lang = language.English.String()
		df.AppScale = 1
		df.Remember.Active = false
		df.DisableBgMusic = false
		df.UsePreRelease = false
	})

	if err := m.syncDataFn(m, true); err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to sync data").Err(err).Msg(log.ErrorManager.String())
	}

	m.handleUIRefresh()
}

func (m *Model) ResetDataAsync() {
	m.commandChan <- ResetDataCommand{}
}
