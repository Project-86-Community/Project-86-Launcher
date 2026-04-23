package p86l

import (
	"context"
	"fmt"
	"os"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"
)

type DownloadManagerModel struct {
	dl fs.Downloader
	ex fs.Extractor
}

func (m *DownloadManagerModel) InstalledVersions(fake bool, afs afero.Fs) ([]Version, error) {
	if fake {
		return []Version{
			{Tag: "v0.0.0-alpha", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
			{Tag: "v1.0.0-beta", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
		}, nil
	}

	versionsDir, err := fs.RealPath(afs, "versions")
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	// All possible executable types across platforms.
	// Platform is determined by the executable file extension/pattern found.
	type exeCandidate struct {
		path string
		os   string // "windows", "darwin", "linux"
	}

	allCandidates := []exeCandidate{
		// Windows
		{"Project-86.exe", "windows"},
		{"Project86.exe", "windows"},
		// Linux
		{"Project-86.x86_64", "linux"},
		{"Project86.x86_64", "linux"},
		{"Project-86.x64", "linux"},
		{"Project86.x64", "linux"},
		{"Project-86", "linux"},
		{"Project86", "linux"},
	}

	var versions []Version
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		tag := entry.Name()
		tagDir := filepath.Join(versionsDir, tag)

		subEntries, err := os.ReadDir(tagDir)
		if err != nil {
			continue
		}

		for _, sub := range subEntries {
			if !sub.IsDir() {
				continue
			}
			gameDir := filepath.Join(tagDir, sub.Name())

			// First, check for macOS .app bundles (they're directories, not files)
			subSubEntries, _ := os.ReadDir(gameDir)
			for _, sse := range subSubEntries {
				if sse.IsDir() && strings.HasSuffix(sse.Name(), ".app") {
					appDir := filepath.Join(gameDir, sse.Name())
					macOSDir := filepath.Join(appDir, "Contents", "MacOS")
					macContents, _ := os.ReadDir(macOSDir)
					for _, mc := range macContents {
						if !mc.IsDir() {
							exePath := filepath.Join(macOSDir, mc.Name())
							runnable := runtime.GOOS == "darwin"
							versions = append(versions, Version{
								Tag:        tag,
								Executable: exePath,
								Runnable:   runnable,
								OS:         "darwin",
							})
							goto nextVersionTag
						}
					}
				}
			}

			// Check for regular executables (Windows/Linux)
			for _, cand := range allCandidates {
				candidate := filepath.Join(gameDir, cand.path)
				info, err := os.Stat(candidate)
				if err == nil && !info.IsDir() {
					var runnable bool
					switch cand.os {
					case "windows":
						runnable = runtime.GOOS == "windows"
					case "linux":
						runnable = runtime.GOOS == "linux"
					}
					versions = append(versions, Version{
						Tag:        tag,
						Executable: candidate,
						Runnable:   runnable,
						OS:         cand.os,
					})
					break
				}
			}
		}

	nextVersionTag:
	}

	return versions, nil
}

func (m *DownloadManagerModel) InstallVersion(ctx context.Context, afs afero.Fs, url string, onProgress func(InstallProgress)) error {
	version, err := fs.ParseVersion(url)
	if err != nil {
		return fmt.Errorf("InstallVersion: could not parse version from URL: %w", err)
	}
	logger.Info.Printf("install started - version: %s url: %s", version, url)

	p := InstallProgress{Phase: types.PhaseDownload}

	zipPath, err := m.dl.Download(ctx, afs, fs.DownloadOptions{
		URL: url,
		Progress: func(done, total int64) {
			p.DownloadDone = done
			p.DownloadTotal = total
			if onProgress != nil {
				onProgress(p)
			}
		},
	})
	if err != nil {
		logger.Error.Printf("download failed [%s]: %v", version, err)
		return fmt.Errorf("InstallVersion: download: %w", err)
	}
	logger.Info.Printf("download complete [%s]: %s", version, zipPath)

	p.Phase = types.PhaseExtract
	if onProgress != nil {
		onProgress(p)
	}
	logger.Info.Printf("extraction started [%s]", version)

	err = m.ex.Extract(ctx, afs, zipPath, fs.ExtractOptions{
		OnFile: func(extracted, total int) {
			p.ExtractDone = extracted
			p.ExtractTotal = total
			if onProgress != nil {
				onProgress(p)
			}
		},
	})
	if err != nil {
		logger.Error.Printf("extraction failed [%s]: %v", version, err)
		return fmt.Errorf("InstallVersion: extract: %w", err)
	}
	logger.Info.Printf("extraction complete [%s]", version)
	logger.Info.Printf("install finished [%s]", version)

	p.Phase = types.PhaseDone
	if onProgress != nil {
		onProgress(p)
	}
	return nil
}
