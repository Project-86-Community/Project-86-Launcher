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
	"strings"

	"github.com/hajimehoshi/guigui"
)

func RemoveLineBreaks(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return input
}

func WrapText(context *guigui.Context, input string, maxWidth int) string {
	charWidthDivisor := 6.2 * context.AppScale()
	charCount := int(float64(maxWidth) / charWidthDivisor)
	input = strings.ReplaceAll(input, "\r\n", "\n")

	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		// Preserve empty lines
		if len(line) == 0 {
			result = append(result, "")
			continue
		}

		// Check if line is already short enough
		if len(line) <= charCount {
			result = append(result, line)
			continue
		}

		// Preserve indentation
		leadingSpace := ""
		trimmedLine := strings.TrimLeft(line, " \t")
		if len(trimmedLine) < len(line) {
			leadingSpace = line[:len(line)-len(trimmedLine)]
		}

		// Only wrap if needed
		words := strings.Fields(trimmedLine)
		if len(words) == 0 {
			result = append(result, line) // Preserve lines with only whitespace
			continue
		}

		var currentLine string
		if len(leadingSpace) > 0 {
			currentLine = leadingSpace
		}

		for _, word := range words {
			// Handle words longer than the wrap limit
			if len(word) > charCount-len(leadingSpace) {
				if len(currentLine) > len(leadingSpace) {
					result = append(result, currentLine)
					currentLine = leadingSpace
				}

				// Split long word
				for i := 0; i < len(word); i += charCount - len(leadingSpace) {
					end := i + charCount - len(leadingSpace)
					if end > len(word) {
						end = len(word)
					}

					if i == 0 {
						result = append(result, currentLine+word[i:end])
					} else {
						result = append(result, leadingSpace+word[i:end])
					}

					if i+charCount-len(leadingSpace) < len(word) {
						currentLine = leadingSpace
					} else {
						currentLine = ""
					}
				}
			} else {
				// Normal word handling
				if len(currentLine) == len(leadingSpace) {
					currentLine += word
				} else if len(currentLine)+1+len(word) <= charCount {
					currentLine += " " + word
				} else {
					result = append(result, currentLine)
					currentLine = leadingSpace + word
				}
			}
		}

		if len(currentLine) > 0 {
			result = append(result, currentLine)
		}
	}

	return strings.Join(result, "\n")
}
