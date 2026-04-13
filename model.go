package p86l

import (
	"context"
	"fmt"
	"os"
	"p86l/assets"
	"p86l/configs"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
)

type InstallProgress struct {
	// Download phase
	DownloadDone  int64
	DownloadTotal int64
	// Extract phase, non-zero means download finished
	ExtractDone  int
	ExtractTotal int

	Phase types.Phase
}

type Version struct {
	Tag        string
	Executable string // absolute path
	Runnable   bool   // whether this version can run on the current OS
	OS         string // "windows", "darwin", "linux"
}

type Model struct {
	fake bool

	mode            types.Mode
	sidebarPosition types.SidebarPosition
	listPosition    types.ListPosition
	t               assets.T

	fs afero.Fs
	dl fs.Downloader
	ex fs.Extractor
}

func NewModel(afs afero.Fs) Model {
	return Model{
		sidebarPosition: types.SidebarPositionRight,
		listPosition:    types.ListPositionTop,
		fs:              afs,
		dl:              fs.GrabDownloader{},
		ex:              fs.FastExtractor{},
	}
}

// Only for testing.

const FakeDownloadURL = "https://github.com/Taliayaya/Project-86/releases/download/v0.0.0-alpha/Project86-v0.0.0-alpha.zip"

func (m *Model) Fake() bool {
	return m.fake
}

func (m *Model) UseFakes() {
	m.dl = fs.FakeDownloader{}
	m.ex = fs.FakeExtractor{}
	m.fake = true
}

func (m *Model) Mode() types.Mode {
	if m.mode == types.ModeUnknown {
		return types.ModeHome
	}
	return m.mode
}

func (m *Model) SetMode(mode types.Mode) {
	m.mode = mode
}

func (m *Model) SetSidebarPosition(pos types.SidebarPosition) {
	m.sidebarPosition = pos
}

func (m *Model) SidebarPosition() types.SidebarPosition {
	return m.sidebarPosition
}

func (m *Model) SetListPosition(pos types.ListPosition) {
	m.listPosition = pos
}

func (m *Model) ListPosition() types.ListPosition {
	return m.listPosition
}

func (m *Model) T() assets.T {
	return m.t
}

func (m *Model) SetT(lang string) {
	m.t = assets.NewT(lang)
}

func (m *Model) OpenFolder(folder types.Folder) {
	if m.fake {
		return
	}

	var rel string
	switch folder {
	case types.FolderRoot:
		rel = "."
	case types.FolderVersions:
		rel = "versions"
	case types.FolderLogs:
		rel = "logs"
	default:
		return
	}

	path, err := fs.RealPath(m.fs, rel)
	if err != nil {
		logger.Warn.Printf("OpenFolder: could not resolve %q: %v", rel, err)
		return
	}

	logger.Info.Printf("opening folder in file manager: %s", rel)
	_ = open.Start(path)
}

func (m *Model) ReadLatestLog() (string, error) {
	if m.fake {
		return "[fake mode - no real log file]", nil
	}

	logsDir, err := fs.RealPath(m.fs, "logs")
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

func (m *Model) InstallVersion(ctx context.Context, url string, onProgress func(InstallProgress)) error {
	version, err := fs.ParseVersion(url)
	if err != nil {
		return fmt.Errorf("InstallVersion: could not parse version from URL: %w", err)
	}
	logger.Info.Printf("install started - version: %s url: %s", version, url)

	p := InstallProgress{Phase: types.PhaseDownload}

	zipPath, err := m.dl.Download(ctx, m.fs, fs.DownloadOptions{
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

	err = m.ex.Extract(ctx, m.fs, zipPath, fs.ExtractOptions{
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

func (m *Model) InstalledVersions() ([]Version, error) {
	if m.fake {
		return []Version{
			{Tag: "v0.0.0-alpha", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
			{Tag: "v1.0.0-beta", Executable: "/fake/Project-86", Runnable: true, OS: "linux"},
		}, nil
	}

	versionsDir, err := fs.RealPath(m.fs, "versions")
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

func (m *Model) OpenVersionFolder(tag string) {
	path, err := fs.RealPath(m.fs, filepath.Join("versions", tag))
	if err != nil {
		return
	}
	logger.Info.Printf("opening version folder: %s", tag)
	_ = open.Start(path)
}

func (m *Model) DeleteVersion(tag string) error {
	path, err := fs.RealPath(m.fs, filepath.Join("versions", tag))
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func (m *Model) CreateShortcut(ver Version) error {
	// Platform-specific shortcut creation.
	switch runtime.GOOS {
	case "windows":
		return nil
	case "darwin":
		return nil
	default:
		return nil
	}
}
