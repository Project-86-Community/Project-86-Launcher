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

package file_test

import (
	"p86l/internal/file"
	"testing"
)

func setup(t *testing.T) *file.Filesystem {
	fs, err := file.NewFilesystem("test")
	if err != nil {
		t.Fatalf("%v", err)
	}
	return fs
}

func TestSaveFile(t *testing.T) {
	fs := setup(t)
	defer func() {
		err := fs.Close()
		if err != nil {
			t.Fatalf("Failed to close fs: %v", err)
		}
	}()

	err := fs.Save("test.txt", []byte(string("test")))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExistFile(t *testing.T) {
	fs := setup(t)
	defer func() {
		err := fs.Close()
		if err != nil {
			t.Fatalf("Failed to close fs: %v", err)
		}
	}()

	value := fs.Exist("test.txt")
	if !value {
		t.Fatal("missing test.txt")
	}
}

func TestLoadFile(t *testing.T) {
	fs := setup(t)
	defer func() {
		err := fs.Close()
		if err != nil {
			t.Fatalf("Failed to close fs: %v", err)
		}
	}()

	_, err := fs.Load("test.txt")
	if err != nil {
		t.Fatalf("%v", err)
	}
}
