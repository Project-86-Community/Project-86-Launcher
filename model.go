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

type Model struct {
	fake bool

	mode types.Mode
	t    assets.T

	fs afero.Fs
	dl fs.Downloader
	ex fs.Extractor
}

func NewModel(afs afero.Fs) Model {
	return Model{
		fs: afs,
		dl: fs.GrabDownloader{},
		ex: fs.FastExtractor{},
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
	m.fs = fs.NewMem()
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
