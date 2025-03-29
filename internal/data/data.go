/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Ilan Mayeux
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Ilan Mayeux
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
	"eightysix/configs"
	"fmt"
	"strconv"

	"github.com/hajimehoshi/guigui"
	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog/log"
)

type Data struct {
	ColorMode guigui.ColorMode
	AppScale  int

	GDataM *gdata.Manager
}

func (d *Data) saveColorMode() error {
	if err := d.GDataM.SaveObjectProp(configs.Data, configs.ColorModeFile, []byte(fmt.Sprintf("%v", d.ColorMode))); err != nil {
		return err
	}
	return nil
}

func (d *Data) saveAppScale() error {
	if err := d.GDataM.SaveObjectProp(configs.Data, configs.AppScaleFile, []byte(fmt.Sprintf("%v", d.AppScale))); err != nil {
		return err
	}
	return nil
}

func (d *Data) InitColorMode() error {
	if d.GDataM.ObjectPropExists(configs.Data, configs.ColorModeFile) {
		colorModeByte, err := d.GDataM.LoadObjectProp(configs.Data, configs.ColorModeFile)
		if err != nil {
			return err
		}
		colorModeData, err := strconv.Atoi(string(colorModeByte))
		if err != nil {
			return err
		}
		d.ColorMode = guigui.ColorMode(colorModeData)

	}
	err := d.saveColorMode()
	return err
}

func (d *Data) InitAppScale() error {
	if d.GDataM.ObjectPropExists(configs.Data, configs.AppScaleFile) {
		appScaleByte, err := d.GDataM.LoadObjectProp(configs.Data, configs.AppScaleFile)
		if err != nil {
			return err
		}
		appScaleData, err := strconv.Atoi(string(appScaleByte))
		if err != nil {
			return err
		}
		d.AppScale = appScaleData
	}
	err := d.saveAppScale()
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

func (d *Data) UpdateData(context *guigui.Context) error {
	if d.ColorMode != context.ColorMode() {
		context.SetColorMode(d.ColorMode)
		if err := d.saveColorMode(); err != nil {
			return err
		}
		log.Info().Int("ColorMode", int(d.ColorMode)).Msg("ColorMode changed")
	}
	if d.AppScale != d.GetAppScale(context.AppScale()) {
		d.SetAppScale(context)
		if err := d.saveAppScale(); err != nil {
			return err
		}
		log.Info().Int("AppScale", d.AppScale).Msg("AppScale changed")
	}
	return nil
}

func (d *Data) HandleDataReset() error {
	d.ColorMode = guigui.ColorModeLight
	d.AppScale = 2

	if err := d.saveColorMode(); err != nil {
		return err
	}
	if err := d.saveAppScale(); err != nil {
		return err
	}
	return nil
}
