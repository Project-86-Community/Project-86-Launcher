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

package p86l

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"p86l/configs"
	"p86l/internal/github"
	"p86l/internal/log"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/dustin/go-humanize"
)

type PlayType int

const (
	PlayInstall PlayType = iota
	PlayUpdate
	PlayPlay
)

func (m *Model) downloadGame(gamePath, gameTag string, asset *github.ReleaseAsset) error {
	client := grab.NewClient()
	req, _ := grab.NewRequest(gamePath, asset.BrowserDownloadURL)

	m.ProgressText(T("model_play.start"))

	resp := client.Do(req)
	m.ProgressText(resp.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for done := false; !done; {
		select {
		case <-t.C:
			m.ProgressText(fmt.Sprintf(
				"%s %s\n\n(%s/%s), %s, %s/s",
				T("model_play.download"),
				gameTag,
				humanize.Bytes(uint64(resp.BytesComplete())),
				humanize.Bytes(uint64(resp.Size())),
				humanize.RelTime(time.Now(), resp.ETA(), "remaining", "ago"),
				humanize.Bytes(uint64(resp.BytesPerSecond())),
			))
		case <-resp.Done:
			done = true
			m.ProgressText(T("model_play.finished"))
		}
	}

	return resp.Err()
}

func (m *Model) unzipGame(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip reader: %w", err)
	}
	defer func() { _ = r.Close() }()

	totalFiles := len(r.File)

	// Extract each files.
	for i, f := range r.File {
		progressInt := int(float64(i+1) / float64(totalFiles) * 100)
		m.ProgressText(fmt.Sprintf("%s %d%%", T("model_play.components"), progressInt))

		err := zipExtractFile(m.fs.Root(), dest, f)
		if err != nil {
			return fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
	}

	return nil
}

func zipExtractFile(fs *os.Root, dest string, f *zip.File) error {
	relPath := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		return zipMkdirAll(fs, relPath, 0755)
	}

	// Create parent directories.
	if err := zipMkdirAll(fs, filepath.Dir(relPath), 0755); err != nil {
		return err
	}

	mode := f.Mode()
	if !mode.IsRegular() {
		mode = 0644
	}

	// Create file.
	outFile, err := fs.OpenFile(relPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	// Open from archive.
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	_, err = io.Copy(outFile, rc)
	return err
}

func zipMkdirAll(root *os.Root, path string, perm os.FileMode) error {
	// Always use a clean directory permission
	perm = perm & os.ModePerm // Strip any non-permission bits
	if perm == 0 {
		perm = 0755
	}

	err := root.Mkdir(path, perm)
	if err == nil || os.IsExist(err) {
		return nil
	}

	// Check if it's an "unsupported file mode" error, try with default perm
	if err != nil && !os.IsNotExist(err) && !os.IsExist(err) {
		err = root.Mkdir(path, 0755)
		if err == nil || os.IsExist(err) {
			return nil
		}
	}

	if parent := filepath.Dir(path); parent != "." && parent != "/" {
		if err := zipMkdirAll(root, parent, 0755); err != nil {
			return err
		}
	}

	err = root.Mkdir(path, 0755) // Always use safe default
	if os.IsExist(err) {
		return nil
	}
	return err
}

func (m *Model) installOrUpdate(isUpdate bool) {
	dataFile := m.Data().Get()
	cacheFile := m.Cache().Get()

	var downloadRelease *github.RepositoryRelease
	var gamePath, gameTag, zipPath string
	var downloadAsset *github.ReleaseAsset
	var resumeVersion string

	if cacheFile.Releases == nil {
		err := T("model_play.missing_releases")
		m.ProgressText(err)
		m.logger.Warn().Str(log.Lifecycle, strings.ToLower(err)).Msg(log.NetworkManager.String())
		return
	}

	if dataFile.UsePreRelease {
		downloadRelease = cacheFile.Releases.PreRelease
		_, debugGameAsset := GetAssets(downloadRelease.Assets)
		downloadAsset = debugGameAsset
		resumeVersion = dataFile.PreReleaseVersion

		gamePath = filepath.Join(configs.FolderBuild, configs.FolderPreRelease)
	} else {
		downloadRelease = cacheFile.Releases.Stable
		gameAsset, _ := GetAssets(downloadRelease.Assets)
		downloadAsset = gameAsset
		resumeVersion = dataFile.GameVersion

		gamePath = filepath.Join(configs.FolderBuild, configs.FolderGame)
	}

	gameTag = downloadRelease.TagName

	if isUpdate {
		gamePath = configs.FolderBuild
		if dataFile.UsePreRelease {
			zipPath = filepath.Join(gamePath, configs.FileUpdatePreRelease)
		} else {
			zipPath = filepath.Join(gamePath, configs.FileUpdateGame)
		}
	} else {
		zipPath = filepath.Join(gamePath, configs.FileBuild)
	}

	// Will delete the game file that's partially downloaded, if a newer version of game came out.
	// Issues are practically rare here, since GUI will not allow this to be executed after Install is done.
	if resumeVersion != "" {
		isNew, err := IsNewVersion(resumeVersion, downloadRelease.TagName)
		if err != nil {
			mErr := T("model_play.unknown_version")
			m.ProgressText(fmt.Sprintf("%s: %v", mErr, err))
			m.logger.Warn().Str(log.Lifecycle, strings.ToLower(mErr)).Err(err).Msg(log.ErrorManager.String())
			return
		}

		if isNew && m.fs.Exist(zipPath) {
			if err := m.fs.Remove(zipPath); err != nil {
				mErr := T("model_play.fail_resume")
				m.ProgressText(fmt.Sprintf("%s: %v", mErr, err))
				m.logger.Warn().Str(log.Lifecycle, strings.ToLower(mErr)).Err(err).Msg(log.ErrorManager.String())
				return
			}
		}
	}

	if downloadAsset == nil {
		err := T("model_play.missing_asset")
		m.ProgressText(err)
		m.logger.Warn().Str(log.Lifecycle, strings.ToLower(err)).Msg(log.NetworkManager.String())
		return
	}
	m.Data().Update(func(df *DataFile) {
		if dataFile.UsePreRelease {
			df.PreReleaseVersion = downloadRelease.TagName
		} else {
			df.GameVersion = downloadRelease.TagName
		}
	})

	m.logger.Info().Str(log.Lifecycle, fmt.Sprintf("downloading file to %s", filepath.Join(m.fs.Path(), zipPath))).Msg(log.FileManager.String())
	if err := m.downloadGame(filepath.Join(m.fs.Path(), zipPath), gameTag, downloadAsset); err != nil {
		m.ProgressText(fmt.Sprintf("%s %v", T("model_play.fail_asset"), err))
		m.logger.Warn().
			Str(log.Lifecycle, fmt.Sprintf("failed to download %s", gameTag)).
			Err(err).
			Caller().
			Msg(log.ErrorManager.String())
		return
	}
	time.Sleep(2 * time.Second)

	m.ProgressText(T("model_play.install_unzip"))
	time.Sleep(2 * time.Second)

	if isUpdate {
		var updatePath string
		if dataFile.UsePreRelease {
			gamePath = filepath.Join(configs.FolderBuild, configs.FolderPreRelease)
			updatePath = filepath.Join(configs.FolderBuild, configs.FolderPreRelease)
		} else {
			gamePath = filepath.Join(configs.FolderBuild, configs.FolderGame)
			updatePath = filepath.Join(configs.FolderBuild, configs.FolderGame)
		}

		if err := m.fs.Root().RemoveAll(updatePath); err != nil {
			mErr := T("model_play.fail_update")
			m.ProgressText(fmt.Sprintf("%s: %v", mErr, err))
			m.logger.Warn().Str(log.Lifecycle, strings.ToLower(mErr)).Err(err).Caller().Msg(log.ErrorManager.String())
			return
		}

	}

	m.logger.Info().Str(log.Lifecycle, "unzipping files").Msg(log.FileManager.String())
	if err := m.unzipGame(filepath.Join(m.fs.Path(), zipPath), gamePath); err != nil {
		mErr := T("model_play.fail_unzip")
		m.ProgressText(fmt.Sprintf("%s: %v", mErr, err))
		m.logger.Warn().Str(log.Lifecycle, strings.ToLower(mErr)).Err(err).Caller().Msg(log.ErrorManager.String())
		return
	}

	m.Data().Update(func(df *DataFile) {
		if dataFile.UsePreRelease {
			df.InstalledPreRelease = downloadRelease.TagName
		} else {
			df.InstalledGame = downloadRelease.TagName
		}
	})
	if m.fs.Exist(zipPath) {
		if err := m.fs.Remove(zipPath); err != nil {
			mErr := T("model_play.fail_artifact")
			m.ProgressText(fmt.Sprintf("%s: %v", mErr, err))
			m.logger.Warn().Str(log.Lifecycle, strings.ToLower(mErr)).Err(err).Msg(log.ErrorManager.String())
			return
		}
	}
	m.logger.Info().Str(log.Lifecycle, "game installation done").Msg(log.FileManager.String())
	m.ProgressText(T("model_play.install_finished"))
	time.Sleep(2 * time.Second)
}

func (m *Model) handlePlay() {
	var exePath string
	data := m.Data()
	dataFile := data.Get()

	if dataFile.UsePreRelease {
		exePath = PathGamePreRelease
	} else {
		exePath = PathGameStable
	}

	if ok := m.CheckFilesCached(exePath); !ok {
		// Failed to find exe
		m.logger.Info().Msg("no exe?")
		return
	}

	cmd := exec.Command(filepath.Join(m.fs.Path(), exePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		// Failed to run exe
		m.logger.Info().Msg("can't run exe?")
		return
	}

	startTime := time.Now()
	data.Update(func(df *DataFile) {
		df.LastPlayed = startTime
	})
	m.logger.Info().Int("Launched", cmd.Process.Pid).Msg(log.AppManager.String())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case <-done:
			sessionTime := time.Since(startTime)
			data.Update(func(df *DataFile) {
				df.TotalPlayTime += sessionTime
			})
			m.logger.Info().Str("Exited after", humanize.RelTime(time.Now(), time.Now().Add(sessionTime), "", "")).Msg(log.AppManager.String())
			return
		case <-sigChan:
			cmd.Process.Kill()
			sessionTime := time.Since(startTime)
			data.Update(func(df *DataFile) {
				df.TotalPlayTime += sessionTime
			})
			m.logger.Info().Str("Stopped after", humanize.RelTime(time.Now(), time.Now().Add(sessionTime), "", "")).Msg(log.AppManager.String())
			return
		}
	}
}

func (m *Model) Play(playType PlayType) {
	m.InProgress(true)
	defer m.InProgress(false)

	switch playType {
	case PlayInstall:
		m.installOrUpdate(false)
	case PlayUpdate:
		m.installOrUpdate(true)
	case PlayPlay:
		m.handlePlay()
	default:
		return
	}

	m.ProgressText("")
	m.logger.Info().Str(log.Lifecycle, "Model.Play is finished").Msg(log.NetworkManager.String())
}
