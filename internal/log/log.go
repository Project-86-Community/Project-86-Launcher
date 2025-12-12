/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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

package log

import (
	"errors"
	"fmt"
	"os"
	"p86l/configs"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var ErrLogFileInvalid = errors.New("failed to create log file")

type Manager int

const (
	UnknownManager Manager = iota
	AppManager
	ErrorManager
	FileManager
	NetworkManager
)

func (m Manager) String() string {
	list := []string{"Unknown", "App", "Error", "File", "Network"}
	return list[m] + "Manager"
}

type Model int

const (
	UnknownModel Model = iota
	MainModel
	DataModel
	CacheModel
)

func (m Model) String() string {
	list := []string{"", "Main", "Data", "Cache"}
	return list[m] + "Model"
}

const (
	Lifecycle      = "lifecycle"
	BackgroundLoop = "background loop"
	InitialFetch   = "initial fetch"
	FetchRateLimit = "fetch rate limit"
	FetchReleases  = "fetch releases"

	Starting = "starting"
	Stopped  = "stopped"
)

// -- errors --

var (
	ErrMkdirAllInvalid    = errors.New("failed to create new folder")
	ErrCompanyPathAppData = errors.New("failed to get appdata")
	ErrRootInvalid        = errors.New("failed to open root")

	ErrFileRemove = errors.New("failed to remove file")
	ErrFileLoad   = errors.New("failed to load file")
	ErrFileSave   = errors.New("failed to save file")

	ErrGithubRequestNew      = errors.New("failed to create new request")
	ErrGithubRequestDo       = errors.New("failed to execute request")
	ErrGithubRequestStatus   = errors.New("github api returned status")
	ErrGithubRequestBodyRead = errors.New("reading body failed")
)

func newLogFile(root *os.Root, path string) (*os.File, *os.File, error) {
	timestamp := time.Now()
	filename := fmt.Sprintf("log_%d-%02d-%02d-%d.txt",
		timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Unix())

	mainPath := filepath.Join(path, filename)
	latestPath := filepath.Join(path, "log-latest.txt")

	main, err := root.Create(mainPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrLogFileInvalid, err)
	}

	latest, err := root.Create(latestPath)
	if err != nil {
		main.Close()
		return nil, nil, fmt.Errorf("failed to create latest log: %w", err)
	}

	return main, latest, nil
}

func NewLogger(VERSION string, fs *os.Root) (*zerolog.Logger, []*os.File, bool, bool, error) {
	switch VERSION {
	case "dev":
		var saveLogs, disableFS, disableAPI bool
		var logger zerolog.Logger
		var logFiles []*os.File

		zerolog.SetGlobalLevel(zerolog.Disabled)

		for token := range strings.SplitSeq(os.Getenv("P86L_DEBUG"), ",") {
			switch token {
			case "log":
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case "logfile":
				saveLogs = true
			case "nofs":
				disableFS = true
			case "noapi":
				disableAPI = true
			}
		}

		if zerolog.GlobalLevel() != zerolog.Disabled {
			lcw := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}

			if saveLogs {
				mainLogFile, latestLogFile, err := newLogFile(fs, filepath.Join(configs.AppName, configs.FolderLogs))
				if err != nil {
					return nil, nil, disableFS, false, err
				}
				logFiles = []*os.File{mainLogFile, latestLogFile}

				multiWriter := zerolog.MultiLevelWriter(lcw, mainLogFile, latestLogFile)
				logger = zerolog.New(multiWriter).With().Timestamp().Logger()
			} else {
				logger = zerolog.New(lcw).With().Timestamp().Logger()
			}
		}

		logger.Info().Bool("Debug", true).Msg(AppManager.String())
		return &logger, logFiles, disableFS, disableAPI, nil
	default:
		mainLogFile, latestLogFile, err := newLogFile(fs, filepath.Join(configs.AppName, configs.FolderLogs))
		if err != nil {
			return nil, nil, false, false, err
		}

		multiWriter := zerolog.MultiLevelWriter(os.Stdout, mainLogFile, latestLogFile)
		logger := zerolog.New(multiWriter).With().Timestamp().Logger()
		return &logger, []*os.File{mainLogFile, latestLogFile}, false, false, nil
	}
}
