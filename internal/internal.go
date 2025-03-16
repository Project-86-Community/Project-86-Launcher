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

package internal

import (
	"archive/zip"
	"eightysix/content"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

func IsInternetReachable() bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 204
}

func IsValidGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		!strings.Contains(filename, "dev")
}

func CheckNewerVersion(currentVersion, newVersion string) (bool, error) {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, fmt.Errorf("invalid current version: %w", err)
	}

	newer, err := version.NewVersion(newVersion)
	if err != nil {
		return false, fmt.Errorf("invalid new version: %w", err)
	}

	return newer.GreaterThan(current), nil
}

func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func OpenFileManager(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform, can't open file manager")
	}

	return cmd.Start()
}

func AddLineBreaks(input string, wordsPerLine int) string {
	words := strings.Fields(input)
	var result strings.Builder

	for i, word := range words {
		result.WriteString(word)

		if i < len(words)-1 {
			if (i+1)%wordsPerLine == 0 {
				result.WriteString("\n")
			} else {
				result.WriteString(" ")
			}
		}
	}

	return result.String()
}

func FindExecutable(rootDir, targetFile string) (string, error) {
	var foundPath string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == targetFile {
			foundPath = path
			return filepath.SkipAll
		}

		return nil
	})

	if foundPath != "" {
		return foundPath, nil
	}

	if err != nil && err != filepath.SkipAll {
		return "", err
	}
	return "", fmt.Errorf("file %s not found in %s", targetFile, rootDir)
}

func RunExecutable(exePath string) error {
	cmd := exec.Command(exePath)

	cmd.Dir = filepath.Dir(exePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func DownloadFile(url string, destPath string) error {
	if destPath == "" {
		parts := strings.Split(url, "/")
		destPath = parts[len(parts)-1]
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	size := resp.ContentLength

	counter := &WriteCounter{
		Total:     size,
		StartTime: time.Now(),
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	return err
}

type WriteCounter struct {
	Total       int64
	Downloaded  int64
	LastPercent float64
	StartTime   time.Time
	LastUpdate  time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += int64(n)
	percent := float64(wc.Downloaded) / float64(wc.Total) * 100
	now := time.Now()

	if percent-wc.LastPercent >= 0.1 || now.Sub(wc.LastUpdate) > 500*time.Millisecond {
		wc.LastPercent = percent
		wc.LastUpdate = now

		elapsed := now.Sub(wc.StartTime).Seconds()
		if elapsed > 0 {
			speed := float64(wc.Downloaded) / elapsed
			remaining := float64(wc.Total-wc.Downloaded) / speed
			eta := formatDuration(remaining)

			content.DownloadStatus = percent
			content.DownloadETA = eta
		} else {
			content.DownloadStatus = percent
		}
	}

	return n, nil
}

func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	} else if seconds < 3600 {
		m := int(seconds) / 60
		s := int(seconds) % 60
		return fmt.Sprintf("%dm %ds", m, s)
	} else {
		h := int(seconds) / 3600
		m := (int(seconds) % 3600) / 60
		return fmt.Sprintf("%dh %dm", h, m)
	}
}

func ExtractZip(zipFile, destination string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(destination, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
