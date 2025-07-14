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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	ErrorManager   = "ErrorManager"
	FileManager    = "FileManager"
	NetworkManager = "NetworkManager"
)

type ErrorType string

const (
	UnknownError ErrorType = "unknown"
	AppError     ErrorType = "app"
	FSError      ErrorType = "filesystem"
	DataError    ErrorType = "data"
	CacheError   ErrorType = "cache"
	NetworkError ErrorType = "network"
)

const (
	// App errors (1001-1999)
	ErrUnknown int = iota + 1001
	ErrBrowserOpen
	ErrGameVersionInvalid
	ErrGameNotExist
)

const (
	// Filesystem errors (2001-2999)
	ErrFSOpenFileManagerInvalid int = iota + 2001

	ErrFSDirInvalid
	ErrFSDirNew
	ErrFSDirNotExist
	ErrFSDirRead
	ErrFSDirRename
	ErrFSDirRemove

	ErrFSFileInvalid
	ErrFSNewFileInvalid
	ErrFSFileWrite
	ErrFSFileNotExist

	// -- root --

	ErrFSRootInvalid

	ErrFSRootDirInvalid
	ErrFSRootDirNew

	ErrFSRootFileInvalid
	ErrFSRootFileNew
	ErrFSRootFileNotExist
	ErrFSRootFileRead
	ErrFSRootFileWrite
	ErrFSRootFileRemove
	ErrFSRootFileClose
)

const (
	// Data errors (3001-3999)
	ErrDataLoad int = iota + 3001
	ErrDataSave
	ErrDataReset
	ErrDataLocaleInvalid
)

const (
	// Cache errors (4001-4999)
	ErrCacheLoad int = iota + 4001
	ErrCacheSave
	ErrCacheReset
	ErrCacheInvalid
	ErrCacheRepoInvalid
	ErrCacheBodyInvalid
	ErrCacheURLInvalid
	ErrCacheAssetsInvalid
)

const (
	// // Network errors (5001-5999)
	ErrNetworkRateLimitInvalid int = iota + 5001
	ErrNetworkCacheRequest
	ErrNetworkDownloadRequest
	ErrNetworkStatusNotOk
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

func (e *Error) String() string {
	if e != nil {
		return fmt.Sprintf("Code: ( %d ), Type: ( %s ), Err: ( %s )", e.Code, string(e.Type), e.Err.Error())
	}
	return ""
}

func (e *Error) LogErr(place, who string) {
	log.Error().Int("Code", e.Code).Any("Type", e.Type).Err((e.Err)).Str(place, who).Msg(ErrorManager)
}

func (e *Error) LogErrStack(place, who string) {
	log.Error().Stack().Int("Code", e.Code).Any("Type", e.Type).Err((e.Err)).Str(place, who).Msg(ErrorManager)
}

func (d *Debug) SetToast(err *Error) {
	log.Warn().Stack().Int("Code", err.Code).Any("Type", err.Type).Err(err.Err).Str("Debug", "SetToast").Msg(ErrorManager)
	d.ToastErr = err
}

func (d *Debug) SetPopup(err *Error) {
	log.Warn().Stack().Int("Code", err.Code).Any("Type", err.Type).Err(err.Err).Str("Debug", "SetPopup").Msg(ErrorManager)
	d.PopupErr = err
}
