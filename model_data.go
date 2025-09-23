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
	"encoding/json"
	"p86l/configs"
	"p86l/internal/file"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/text/language"
)

type Pages int

const (
	PageHome Pages = iota
	PagePlay
	PageSettings
	PageAbout
)

type DataData struct {
	Lang           string  `json:"lang"`
	UseDarkmode    bool    `json:"use_darkmode"`
	AppScale2      float64 `json:"app_scale_2"`
	DisableBgMusic bool    `json:"disable_bgm"`
	UsePreRelease  bool    `json:"use_pre_release"`
}

type DataModel struct {
	fs *file.Filesystem

	isNew bool
	data  *DataData

	page Pages

	lang           language.Tag
	useDarkmode    bool
	scale          float64
	disableBgMusic bool
	usePreRelease  bool
}

// -- Getters for DataModel --

func (c *DataModel) Path() string {
	return filepath.Join(configs.AppName, configs.FileData)
}

func (d *DataModel) IsNew() bool {
	return d.isNew
}

func (d *DataModel) Data() error {
	if d.data == nil {
		d.isNew = true
		d.New()
	}
	if err := d.Start(); err != nil {
		return err
	}

	return nil
}

// -- data

func (d *DataModel) Page() Pages {
	return d.page
}

func (d *DataModel) Lang() language.Tag {
	return d.lang
}

func (d *DataModel) UseDarkmode() bool {
	return d.useDarkmode
}

func (d *DataModel) AppScale() float64 {
	return d.scale
}

func (d *DataModel) DisableBgMusic() bool {
	return d.disableBgMusic
}

func (d *DataModel) UsePreRelease() bool {
	return d.usePreRelease
}

func (d *DataModel) Load() error {
	b, err := d.fs.Load(d.Path())
	if err != nil {
		return err
	}

	var _d *DataData
	if err := json.Unmarshal(b, &_d); err != nil {
		return err
	}
	d.data = _d

	return nil
}

// -- Setters for DataModel --

func (d *DataModel) SetFS(fs *file.Filesystem) {
	d.fs = fs
}

func (d *DataModel) SetData(data *DataData) {
	d.data = data
}

func (d *DataModel) SetPage(page Pages) {
	d.page = page
}

func (d *DataModel) SetLang(value language.Tag) {
	d.lang = value
}

func (d *DataModel) SetUseDarkmode(value bool) {
	d.useDarkmode = value
}

func (d *DataModel) SetAppScale(value float64) {
	d.scale = value
}

func (d *DataModel) SetDisableBgMusic(player *audio.Player, value bool) {
	d.disableBgMusic = value
	switch value {
	case true:
		player.Pause()
	case false:
		player.Play()
	}
}

func (d *DataModel) SetUsePreRelease(value bool) {
	d.usePreRelease = value
}

func (d *DataModel) Save(data DataData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = d.fs.Save(d.Path(), b)
	if err != nil {
		return err
	}

	return nil
}

// -- common --

func (d *DataModel) New() {
	d.data = &DataData{
		Lang:           language.English.String(),
		UseDarkmode:    false,
		AppScale2:      1.0,
		DisableBgMusic: false,
		UsePreRelease:  false,
	}
}

func (d *DataModel) Start() error {
	data := d.data

	lang, err := language.Parse(data.Lang)
	if err != nil {
		return err
	}

	d.lang = lang
	d.useDarkmode = data.UseDarkmode
	d.scale = data.AppScale2
	d.disableBgMusic = data.DisableBgMusic
	d.usePreRelease = data.UsePreRelease

	return nil
}
