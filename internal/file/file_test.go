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
	"os"
	"p86l/configs"
	"p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"testing"
	"time"
)

func setup(t *testing.T) (*debug.Debug, *file.AppFS) {
	e := &debug.Debug{}

	companyPath, dErr := file.GetCompanyPath(e, "test")
	if dErr.Err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}
	root, err := os.OpenRoot(companyPath)
	if err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", debug.ErrDirNotFound, string(debug.FSError), err)
	}

	a, dErr := file.NewFS(e, root, companyPath, configs.AppName)
	if dErr.Err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}

	return e, a
}

func TestDirs(t *testing.T) {
	_, fs := setup(t)
	t.Logf("%#v", fs)
}

func TestSaveFiles(t *testing.T) {
	e, fs := setup(t)
	dErr := fs.Save(e, configs.Data, configs.ColorModeFile, []byte("Lena"))
	if dErr.Err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}
}

func TestLoadFiles(t *testing.T) {
	e, fs := setup(t)
	bytes, dErr := fs.Load(e, configs.Data, configs.ColorModeFile)
	if dErr.Err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}
	t.Log(string(bytes))
}

func TestStatAppDir(t *testing.T) {
	e, fs := setup(t)
	if dErr := fs.IsDir(e, filepath.Join(fs.CompanyDirPath, configs.AppName, configs.Data, configs.ColorModeFile)); dErr.Err != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}
}

func TestDataFiles(t *testing.T) {
	e, fs := setup(t)
	time.Sleep(time.Duration(2 * time.Second))
	dErr := fs.ResetData(e)
	if dErr != nil {
		t.Fatalf("Code: %d, Type: %s, Err: %v", dErr.Code, string(dErr.Type), dErr.Err)
	}
}
