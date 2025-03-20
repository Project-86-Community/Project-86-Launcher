/*
 * SPDX-License-Identifier: GPL-3.0-only
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

package internal

import (
	"eightysix/content"
	"fmt"

	"github.com/hajimehoshi/guigui"
)

func SaveDarkMode() error {
	if err := content.Mgdata.SaveObjectProp("cache", "darkmode.data", []byte(fmt.Sprintf("%v", content.ColorMode))); err != nil {
		return err
	}
	return nil
}

func SaveAppScale() error {
	if err := content.Mgdata.SaveObjectProp("cache", "appscale.data", []byte(fmt.Sprintf("%v", content.AppScale))); err != nil {
		return err
	}
	return nil
}

func GetAppScale(scale float64) int {
	var appScale int

	switch scale {
	case 0.5: // 50%
		appScale = 0
	case 0.75: // 75%
		appScale = 1
	case 1.0: // 100%
		appScale = 2
	case 1.25: // 125%
		appScale = 3
	case 1.50: // 150%
		appScale = 4
	}

	return appScale
}

func SetAppScale(context *guigui.Context) {
	switch content.AppScale {
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
