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
	"archive/zip"
	gctx "context"
	"fmt"
	"io"
	"net/http"
	"os"
	"p86l/configs"
	pd "p86l/internal/debug"
	"p86l/internal/file"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/hashicorp/go-version"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

func LoadB(context *guigui.Context, model *Model, loadType string) *pd.Error {
	switch loadType {
	case "data":
		if err := FS.IsDirR(E, FS.FileDataPath()); err != nil {
			log.Info().Str("Data", "data not found, creating data...").Str("utils", "loadB").Msg(pd.FileManager)
			d := NewData()
			d.Log()
			model.data.file = d
			return model.data.Save()
		}
	case "cache":
		if err := FS.IsDirR(E, FS.FileCachePath()); err != nil {
			log.Info().Str("Cache", "cache not found").Str("utils", "loadB").Msg(pd.FileManager)
			return nil
		}
	}

	switch loadType {
	case "data":
		d, err := LoadData()
		if err != nil {
			return err
		}

		tag, rErr := language.Parse(d.Locale)
		if rErr != nil {
			return E.New(rErr, pd.DataError, pd.ErrDataLocaleInvalid)
		}
		model.data.SetLocale(context, tag)
		model.data.SetAppScale(context, d.AppScale)
		model.data.SetColorMode(context, d.ColorMode)
		model.data.SetUsePreRelease(d.UsePreRelease)
		return model.data.SetGameVersion(d.GameVersion)
	case "cache":
		c, err := LoadCache()
		if err != nil {
			return err
		}
		if err := c.Validate(E); err == nil {
			model.cache.valid = true
		}
		model.cache.file = *c
	}
	return nil
}

func SetLanguage(lang string) {
	LLocalizer = i18n.NewLocalizer(LBundle, lang)
}

func T(key string) string {
	lMsg, err := LLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return fmt.Sprintf("!{%s}", key)
	}

	return lMsg
}

func translateGT(body string, target string) string {
	result, err := t.Translate(body, "auto", target)
	if err != nil {
		return "?"
	}
	return result.Text
}

func OpenBrowser(url string) {
	log.Info().Str("Url", url).Msg("OpenBrowser")
	if err := browser.OpenURL(url); err != nil {
		E.SetPopup(E.New(err, pd.AppError, pd.ErrBrowserOpen))
	}
}

func IsValidPreGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip")
}

func IsValidGameFile(filename string) bool {
	return strings.Contains(filename, "Project86-v") &&
		strings.Contains(filename, ".zip") &&
		!strings.Contains(filename, "dev")
}

func CheckNewerVersion(currentVersion, newVersion string) (bool, *pd.Error) {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, E.New(fmt.Errorf("invalid current version: %w", err), pd.AppError, pd.ErrGameVersionInvalid)
	}

	newer, err := version.NewVersion(newVersion)
	if err != nil {
		return false, E.New(fmt.Errorf("invalid new version: %w", err), pd.AppError, pd.ErrGameVersionInvalid)
	}

	return newer.GreaterThan(current), nil
}

// -- downloading --

func GetPreRelease() (*github.RepositoryRelease, error) {
	ctx := gctx.Background()
	opt := &github.ListOptions{
		PerPage: 100,
	}

	for {
		rs, resp, err := GithubClient.Repositories.ListReleases(ctx, configs.RepoOwner, configs.RepoName, opt)
		if err != nil {
			return nil, err
		}

		for _, r := range rs {
			if r.Prerelease != nil && *r.Prerelease {
				return r, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil, fmt.Errorf("no pre-releases found")
}

func DownloadGame(model *Model, filename, src string, preRelease bool) *pd.Error {
	var game string
	if preRelease {
		game = "pregame"
	} else {
		game = "game"
	}
	gameZip := fmt.Sprintf("%s.zip", game)

	model.SetProgress("Downloading...")

	err := DownloadFile(model, filename, src, gameZip)
	if err != nil {
		model.SetProgress("")
		return err
	}

	if FS.IsDirR(E, FS.DirGamePath()) == nil {
		model.SetProgress("Removing old files...")

		rErr := os.RemoveAll(filepath.Join(FS.CompanyDirPath, "build", game))
		if rErr != nil {
			model.SetProgress("")
			return E.New(rErr, pd.FSError, pd.ErrFSDirRemove)
		}
	}

	model.SetProgress("Extracting...")

	err = unzip(filepath.Join(FS.CompanyDirPath, gameZip), filepath.Join(FS.CompanyDirPath, "build", game))
	if err != nil {
		model.SetProgress("")
		return err
	}

	model.SetProgress("Cleaning...")

	rErr := FS.Root.Remove(gameZip)
	if rErr != nil {
		model.SetProgress("")
		return E.New(rErr, pd.FSError, pd.ErrFSRootFileRemove)
	}

	model.SetProgress("")

	return nil
}

func DownloadFile(model *Model, filename, src, dest string) *pd.Error {
	var resumePos int64 = 0
	if info, err := FS.Root.Stat(dest); err == nil {
		resumePos = info.Size()
	}

	if resumePos > 0 {
		// Make a HEAD request to get the file size without downloading.
		headReq, err := http.NewRequest("HEAD", src, nil)
		if err != nil {
			return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
		}

		client := &http.Client{}
		headResp, err := client.Do(headReq)
		if err != nil {
			return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
		}
		err = headResp.Body.Close()
		if err != nil {
			return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
		}

		if headResp.StatusCode == http.StatusOK {
			contentLength := headResp.ContentLength
			// If the file size matches, skip download.
			if contentLength > 0 && resumePos == contentLength {
				// File is already fully downloaded.
				return nil
			}
		}
	}

	var out io.WriteCloser
	var err error
	if resumePos > 0 {
		out, err = FS.Root.OpenFile(dest, os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		out, err = FS.Root.Create(dest)
	}
	if err != nil {
		return E.New(err, pd.FSError, pd.ErrFSRootFileNew)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
		}
	}()

	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
	}

	if resumePos > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", resumePos))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
		}
	}()

	expectedStatus := http.StatusOK
	if resumePos > 0 {
		expectedStatus = http.StatusPartialContent
	}

	if resp.StatusCode != expectedStatus {
		if resumePos > 0 && resp.StatusCode != http.StatusPartialContent {
			err := out.Close()
			if err != nil {
				return E.New(err, pd.FSError, pd.ErrFSRootFileClose)
			}
			return DownloadFileFromScratch(model, filename, src, dest)
		}
		return E.New(fmt.Errorf("bad status: %s", resp.Status), pd.NetworkError, pd.ErrNetworkStatusNotOk)
	}

	var totalSize int64
	if resumePos > 0 {
		if contentRange := resp.Header.Get("Content-Range"); contentRange != "" {
			if parts := strings.Split(contentRange, "/"); len(parts) == 2 {
				if size, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					totalSize = size
				}
			}
		}
	} else {
		totalSize = resp.ContentLength
	}

	p := &ProgressTracker{
		model:       model,
		filename:    filename,
		totalSize:   totalSize,
		currentSize: resumePos,
		startTime:   time.Now(),
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, p))
	if err != nil {
		return E.New(err, pd.FSError, pd.ErrFSRootFileWrite)
	}

	return nil
}

func DownloadFileFromScratch(model *Model, filename, src, dest string) *pd.Error {
	// Remove existing file
	err := FS.Root.Remove(dest)
	if err != nil {
		return E.New(err, pd.FSError, pd.ErrFSRootFileRemove)
	}

	// Create new file
	out, err := FS.Root.Create(dest)
	if err != nil {
		return E.New(err, pd.FSError, pd.ErrFSRootFileNew)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
		}
	}()

	resp, err := http.Get(src)
	if err != nil {
		return E.New(err, pd.NetworkError, pd.ErrNetworkDownloadRequest)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return E.New(fmt.Errorf("bad status: %s", resp.Status), pd.NetworkError, pd.ErrNetworkStatusNotOk)
	}

	p := &ProgressTracker{
		model:       model,
		filename:    filename,
		totalSize:   resp.ContentLength,
		currentSize: 0,
		startTime:   time.Now(),
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, p))
	if err != nil {
		err := FS.Root.Remove(dest)
		if err != nil {
			return E.New(err, pd.FSError, pd.ErrFSRootFileRemove)
		}
		return E.New(err, pd.FSError, pd.ErrFSRootFileWrite)
	}

	return nil
}

func unzip(src, dest string) *pd.Error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return E.New(err, pd.FSError, pd.ErrFSRootFileRead)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
		}
	}()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return E.New(err, pd.FSError, pd.ErrFSDirNew)
	}

	// Find common root directory if it exists
	var rootDir string
	if len(r.File) > 0 {
		firstPath := r.File[0].Name
		rootDir = strings.Split(firstPath, "/")[0] + "/"

		// Verify all files have this common root
		for _, f := range r.File {
			if !strings.HasPrefix(f.Name, rootDir) {
				rootDir = "" // No common root
				break
			}
		}
	}

	extractAndWriteFile := func(f *zip.File) *pd.Error {
		// Skip the root directory itself
		if f.Name == rootDir || f.Name == strings.TrimSuffix(rootDir, "/") {
			return nil
		}

		rc, err := f.Open()
		if err != nil {
			return E.New(err, pd.FSError, pd.ErrFSFileInvalid)
		}
		defer func() {
			err := rc.Close()
			if err != nil {
				E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
			}
		}()

		// Strip the root directory from the path
		relPath := f.Name
		if rootDir != "" {
			relPath = strings.TrimPrefix(f.Name, rootDir)
		}

		path := filepath.Join(dest, relPath)

		// Check for ZipSlip vulnerability.
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return E.New(fmt.Errorf("illegal file path: %s", path), pd.FSError, pd.ErrFSFileInvalid)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return E.New(err, pd.FSError, pd.ErrFSDirNew)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return E.New(err, pd.FSError, pd.ErrFSDirNew)
			}

			fi, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return E.New(err, pd.FSError, pd.ErrFSFileInvalid)
			}
			defer func() {
				err := fi.Close()
				if err != nil {
					E.SetToast(E.New(err, pd.FSError, pd.ErrFSRootFileClose))
				}
			}()

			_, err = io.Copy(fi, rc)
			if err != nil {
				return E.New(err, pd.FSError, pd.ErrFSFileWrite)
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// -- Funcs for loading and saving --

func LoadData() (*file.Data, *pd.Error) {
	b, err := FS.Load(E, FS.FileDataPath())
	if err != nil {
		return nil, err
	}

	d, err := FS.DecodeData(E, b)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func LoadCache() (*file.Cache, *pd.Error) {
	b, err := FS.Load(E, FS.FileCachePath())
	if err != nil {
		return nil, err
	}

	c, err := FS.DecodeCache(E, b)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func SaveData(d file.Data) *pd.Error {
	b, err := FS.EncodeData(E, d)
	if err != nil {
		return err
	}

	err = FS.Save(E, FS.FileDataPath(), b)
	if err != nil {
		return err
	}

	return nil
}

func SaveCache(c file.Cache) *pd.Error {
	b, err := FS.EncodeCache(E, c)
	if err != nil {
		return err
	}

	err = FS.Save(E, FS.FileCachePath(), b)
	if err != nil {
		return err
	}

	return nil
}
