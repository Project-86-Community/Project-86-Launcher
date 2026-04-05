package p86l

import (
	"context"
	"log"
	"p86l/assets"
	"p86l/internal/fs"
	"p86l/internal/types"

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
		return
	}
	_ = open.Start(path)
	log.Println(path)
}

func (m *Model) InstallVersion(ctx context.Context, url string, onProgress func(InstallProgress)) error {
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
		return err
	}

	p.Phase = types.PhaseExtract
	if onProgress != nil {
		onProgress(p)
	}

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
		return err
	}

	p.Phase = types.PhaseDone
	if onProgress != nil {
		onProgress(p)
	}
	return nil
}
