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

package open_test

import (
	"p86l/internal/open"
	"testing"
)

func TestOpen(t *testing.T) {
	input := "https://google.com/"
	err := open.Open(input)
	if err != nil {
		t.Errorf("open.Start(\"%s\") threw an error: %s", input, err)
	}

	input = "xxxxxxxxxxxxxxx"
	err = open.Open(input)
	if err != nil {
		t.Errorf("open.Start(\"%s\") shouldn't even fail on invalid input: %s", input, err)
	}
}
