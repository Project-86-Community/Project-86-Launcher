// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package service

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"p86l/configs"
	"p86l/internal/fs"

	"github.com/spf13/afero"
)

type logService struct {
	afs afero.Fs
}

func NewLogService(afs afero.Fs) LogService {
	return &logService{afs: afs}
}

func (s *logService) ReadLatestLog() (string, error) {
	logsDir, err := fs.RealPath(s.afs, "logs")
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return "", err
	}

	var logs []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, configs.LogPrefix) && strings.HasSuffix(name, configs.LogExt) {
			logs = append(logs, name)
		}
	}

	if len(logs) == 0 {
		return "[no previous session log found]", nil
	}

	sort.Strings(logs)
	prev := filepath.Join(logsDir, logs[len(logs)-1])

	data, err := os.ReadFile(prev)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
