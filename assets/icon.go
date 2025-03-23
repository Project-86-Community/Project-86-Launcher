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

package assets

import (
	"bytes"
	_ "embed"
	"image"

	ico "github.com/biessek/golang-ico"
)

//go:embed eightysix.ico
var EightySixIcon []byte

// GetIconImages loads the embedded icon file and returns all images inside it
func GetIconImages() ([]image.Image, error) {
	var IconImages []image.Image

	// Create a reader from the embedded bytes
	reader := bytes.NewReader(EightySixIcon)

	// Decode the ICO file
	icons, err := ico.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	// Add all images from the ICO to our slice
	IconImages = append(IconImages, icons...)

	return IconImages, nil
}
