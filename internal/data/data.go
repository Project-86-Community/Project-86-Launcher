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

package data

import (
	"fmt"
	"p86l/configs"
	"p86l/internal/debug"
	"strconv"

	"github.com/hajimehoshi/guigui"
	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog/log"
)

type Data struct {
	GDataM *gdata.Manager

	ColorMode guigui.ColorMode
	AppScale  int
}

func (d *Data) saveColorMode(appDebug *debug.Debug) *debug.Error {
	if err := d.GDataM.SaveObjectProp(configs.Data, configs.ColorModeFile, []byte(fmt.Sprintf("%v", d.ColorMode))); err != nil {
		return appDebug.New(err, debug.DataError, debug.ErrColorModeSave)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *Data) saveAppScale(appDebug *debug.Debug) *debug.Error {
	if err := d.GDataM.SaveObjectProp(configs.Data, configs.AppScaleFile, []byte(fmt.Sprintf("%v", d.AppScale))); err != nil {
		return appDebug.New(err, debug.DataError, debug.ErrAppScaleSave)
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *Data) InitColorMode(appDebug *debug.Debug) *debug.Error {
	if d.GDataM.ObjectPropExists(configs.Data, configs.ColorModeFile) {
		colorModeByte, err := d.GDataM.LoadObjectProp(configs.Data, configs.ColorModeFile)
		if err != nil {
			return appDebug.New(err, debug.DataError, debug.ErrColorModeLoad)
		}
		colorModeData, err := strconv.Atoi(string(colorModeByte))
		if err != nil {
			return appDebug.New(err, debug.DataError, debug.ErrColorModeLoad)
		}
		d.ColorMode = guigui.ColorMode(colorModeData)
	}
	err := d.saveColorMode(appDebug)
	return err
}

func (d *Data) InitAppScale(appDebug *debug.Debug) *debug.Error {
	if d.GDataM.ObjectPropExists(configs.Data, configs.AppScaleFile) {
		appScaleByte, err := d.GDataM.LoadObjectProp(configs.Data, configs.AppScaleFile)
		if err != nil {
			return appDebug.New(err, debug.DataError, debug.ErrAppScaleLoad)
		}
		appScaleData, err := strconv.Atoi(string(appScaleByte))
		if err != nil {
			return appDebug.New(err, debug.DataError, debug.ErrAppScaleLoad)
		}
		d.AppScale = appScaleData
	}
	err := d.saveAppScale(appDebug)
	return err
}

func (d *Data) GetAppScale(scale float64) int {
	switch scale {
	case 0.5: // 50%
		return 0
	case 0.75: // 75%
		return 1
	case 1.0: // 100%
		return 2
	case 1.25: // 125%
		return 3
	case 1.50: // 150%
		return 4
	}

	return -1
}

func (d *Data) SetAppScale(context *guigui.Context) {
	switch d.AppScale {
	case 0: // 50%
		context.SetAppScale(0.5)
	case 1: // 75%
		context.SetAppScale(0.75)
	case 2: // 100%
		context.SetAppScale(1)
	case 3: // 125%
		context.SetAppScale(1.25)
	case 4: // 150%
		context.SetAppScale(1.50)
	}
}

func (d *Data) UpdateData(context *guigui.Context, appDebug *debug.Debug) *debug.Error {
	if d.ColorMode != context.ColorMode() {
		context.SetColorMode(d.ColorMode)
		if err := d.saveColorMode(appDebug); err.Err != nil {
			return err
		}
		log.Info().Int("ColorMode", int(d.ColorMode)).Msg("ColorMode changed")
	}
	if d.AppScale != d.GetAppScale(context.AppScale()) {
		d.SetAppScale(context)
		if err := d.saveAppScale(appDebug); err.Err != nil {
			return err
		}
		log.Info().Int("AppScale", d.AppScale).Msg("AppScale changed")
	}
	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown)
}

func (d *Data) HandleDataReset(appDebug *debug.Debug) *debug.Error {
	d.ColorMode = guigui.ColorModeLight
	d.AppScale = 2

	if err := d.saveColorMode(appDebug); err.Err != nil {
		return err
	}
	if err := d.saveAppScale(appDebug); err.Err != nil {
		return err
	}

	return appDebug.New(nil, debug.UnknownError, debug.ErrUnknown, "Handle data reset")
}
