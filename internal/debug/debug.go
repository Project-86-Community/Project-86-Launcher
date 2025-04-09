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
	AppError      ErrorType = "app"
	FSError       ErrorType = "filesystem"
	InternetError ErrorType = "internet"
	DataError     ErrorType = "data"
	CacheError    ErrorType = "cache"
)

const (
	// App errors (1001-1999)
	ErrUnknown int = iota + 1001
	ErrBrowserOpen

	// Filesystem errors (2001-2999)
	ErrGDataOpenFailed int = iota + 2001
	ErrIconNotFound
	ErrDirNotFound
	ErrNewDirFailed
	ErrNewFileFailed
	ErrOpenFolderFailed
	ErrFileNotFound
	ErrFolderClear

	// Data errors (3001-3999)
	ErrColorModeLoad int = iota + 3001
	ErrAppScaleLoad
	ErrColorModeSave
	ErrAppScaleSave
	ErrColorModeClear
	ErrAppScaleClear

	// Cache errors (4001-4999)
	ErrChangelogLoad int = iota + 4001
	ErrChangelogSave
	ErrChangelogClear
	ErrChangelogNetwork
)

type Error struct {
	Err     error
	Type    ErrorType
	Code    int
	Message string
}

type Debug struct {
	ToastErr *Error
	PopupErr *Error
}

func (d *Debug) New(err error, errType ErrorType, code int, message ...string) *Error {
	if err != nil {
		return &Error{
			Err:  errors.New(err.Error()),
			Type: errType,
			Code: code,
		}
	}
	if len(message) > 0 {
		return &Error{
			Err:     nil,
			Type:    errType,
			Code:    code,
			Message: message[0],
		}
	}
	return &Error{
		Err:  nil,
		Type: errType,
		Code: code,
	}
}

func (d *Debug) SetToast(err *Error) {
	log.Error().Stack().Int("Code", err.Code).Str("Type", string(err.Type)).Err(err.Err).Msg("Toast error")
	d.ToastErr = err
}

func (d *Debug) SetPopup(err *Error) {
	log.Error().Stack().Int("Code", err.Code).Str("Type", string(err.Type)).Err(err.Err).Msg("Toast error")
	d.PopupErr = err
}

// const (
// 	// Core launcher errors (1-99)
// 	ErrLauncherInit    = 1
// 	ErrLauncherUpdate  = 2
// 	ErrConfigCorrupted = 3
//
// 	// Authentication errors (100-199)
// 	ErrLoginFailed       = 100
// 	ErrSessionExpired    = 101
// 	ErrAccountLocked     = 102
// 	ErrTwoFactorRequired = 103
//
// 	// Game library errors (200-299)
// 	ErrLibraryCorrupted  = 200
// 	ErrGameNotFound      = 201
// 	ErrGameMetadataFetch = 202
// 	ErrGameArtworkFetch  = 203
//
// 	// Download/installation errors (300-399)
// 	ErrDownloadFailed    = 300
// 	ErrInsufficientSpace = 301
// 	ErrChecksumMismatch  = 302
// 	ErrInstallCorrupted  = 303
// 	ErrPatchFailed       = 304
//
// 	// Game execution errors (400-499)
// 	ErrGameLaunchFailed  = 400
// 	ErrMissingDependency = 401
// 	ErrIncompatibleOS    = 402
// 	ErrInsufficientHW    = 403
//
// 	// Network errors (500-599)
// 	ErrServerUnavailable = 500
// 	ErrConnectionLost    = 501
// 	ErrSlowConnection    = 502
// 	ErrCDNFailure        = 503
//
// 	// User profile errors (600-699)
// 	ErrProfileCorrupted = 600
// 	ErrSaveGameSync     = 601
// 	ErrAchievementSync  = 602
// 	ErrFriendListFetch  = 603
//
// 	// Store/purchase errors (700-799)
// 	ErrPaymentFailed      = 700
// 	ErrPurchaseIncomplete = 701
// 	ErrEntitlementIssue   = 702
// 	ErrStoreFetchFailed   = 703
// )
//
// const (
// 	// General errors (1-99)
// 	ErrUnknown  = 1
// 	ErrInternal = 2
//
// 	// Network errors (100-199)
// 	ErrNetworkUnavailable = 100
// 	ErrTimeoutExceeded    = 101
// 	ErrBadResponse        = 102
//
// 	// Database errors (200-299)
// 	ErrDBConnection    = 200
// 	ErrQueryFailed     = 201
// 	ErrRecordNotFound  = 202
// 	ErrDuplicateRecord = 203
//
// 	// Validation errors (300-399)
// 	ErrInvalidInput  = 300
// 	ErrMissingField  = 301
// 	ErrInvalidFormat = 302
//
// 	// Auth errors (400-499)
// 	ErrUnauthorized = 400
// 	ErrForbidden    = 401
// 	ErrTokenExpired = 402
//
// 	// Filesystem errors (500-599)
// 	ErrFileNotFound     = 500
// 	ErrPermissionDenied = 501
// 	ErrDiskFull         = 502
//
// 	// Cache errors (600-699)
// 	ErrCacheMiss    = 600
// 	ErrCacheExpired = 601
// )
