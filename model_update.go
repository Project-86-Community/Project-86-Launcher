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
	"p86l/internal/log"

	"github.com/pkg/browser"
	"github.com/rs/zerolog"
)

type UpdateModel struct {
	urlChan chan string
}

func (u *UpdateModel) OpenURL(url string) {
	select {
	case u.urlChan <- url:
	default:
	}
}

func (u *UpdateModel) Run(logger *zerolog.Logger) error {
	go func() {
		for openUrl := range u.urlChan {
			if err := browser.OpenURL(openUrl); err != nil {
				logger.Warn().Str("openURL", openUrl).Msg(log.ErrorManager.String())
			}
		}
	}()

	return nil
}
