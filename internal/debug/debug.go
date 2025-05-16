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

package debug

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ErrorType string

const (
	UnknownError  ErrorType = "unknown"
	FSError       ErrorType = "filesystem"
	DataError     ErrorType = "data"
	CacheError    ErrorType = "cache"
	InternetError ErrorType = "internet"
)

const (
	// App errors (1001-1999)
	ErrUnknown int = iota + 1001
	ErrBrowserOpen
)

const (
	// Filesystem errors (2001-2999)
	ErrFSFileInvalid int = iota + 2001
	ErrFSDirInvalid
	ErrFSDirNotExist
	ErrFSDirNew
	ErrFSNewFileInvalid
	ErrFSFileNotExist
	ErrFSOpenFileManagerInvalid
	ErrFSRootInvalid
	ErrFSRootDirInvalid
	ErrFSRootDirNew
	ErrFSRootFileInvalid
	ErrFSRootFileNew
	ErrFSRootFileNotExist
	ErrFSRootFileRead
	ErrFSRootFileWrite
)

const (
	// Data errors (3001-3999)
	ErrDataLoad int = iota + 3001
	ErrDataSave
	ErrDataReset
)

const (
	// Cache errors (4001-4999)
	ErrCacheLoad int = iota + 4001
	ErrCacheSave
	ErrCacheReset
	ErrCacheRepoInvalid
	ErrCacheBodyInvalid
	ErrCacheURLInvalid
	ErrCacheAssetsInvalid
)

const (
	// // Internet errors (5001-5999)
	ErrInternetRateLimitInvalid int = iota + 5001
	ErrInternetRequestInvalid
	ErrInternetRequestNew
	ErrInternetRequestClose
	ErrInternetCacheInvalid

	// ErrConnection
	// ErrStatusCode
	// ErrNewRequest
	// ErrFailedDownload
)

type Error struct {
	Err  error
	Type ErrorType
	Code int
}

type Debug struct {
	ToastErr *Error
	PopupErr *Error
}

func (d *Debug) New(err error, errType ErrorType, code int) *Error {
	if err != nil {
		return &Error{
			Err:  errors.Wrap(err, "Debug"),
			Type: errType,
			Code: code,
		}
	}
	return &Error{
		Err:  nil,
		Type: errType,
		Code: code,
	}
}

func (d *Debug) SetToast(err *Error) {
	log.Info().Stack().Int("Code", err.Code).Str("Type", string(err.Type)).Err(err.Err).Msg("SetToast")
	d.ToastErr = err
}

func (d *Debug) SetPopup(err *Error) {
	log.Info().Stack().Int("Code", err.Code).Str("Type", string(err.Type)).Err(err.Err).Msg("SetPopup")
	d.PopupErr = err
}
