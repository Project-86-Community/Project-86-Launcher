package p86l

import (
	"p86l/assets"
	"p86l/internal/fs"
	"p86l/internal/logger"
	"p86l/internal/types"

	"github.com/hajimehoshi/ebiten/v2/audio"
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
	fake      bool
	fakeError bool

	fs afero.Fs
	t  *assets.T

	wvCh   chan<- WebviewRequest
	player *audio.Player

	mode         types.Mode
	listPosition types.ListPosition

	sidebar SidebarModel
	dm      DownloadManagerModel
}

func NewModel(afs afero.Fs, wvCh chan<- WebviewRequest, player *audio.Player) *Model {
	return &Model{
		fs:           afs,
		wvCh:         wvCh,
		player:       player,
		listPosition: types.ListPositionTop,
		sidebar:      SidebarModel{sidebarPosition: types.SidebarPositionRight},
		dm:           DownloadManagerModel{dl: fs.GrabDownloader{}, ex: fs.FastExtractor{}},
	}
}

// Only for testing.
const FakeDownloadURL = "https://github.com/Taliayaya/Project-86/releases/download/v0.0.0-alpha/Project86-v0.0.0-alpha.zip"

func (m *Model) Fake() bool {
	return m.fake
}

func (m *Model) FakeError() bool {
	return m.fakeError
}

func (m *Model) UseFakes(fakeError bool) {
	m.dm.dl = fs.FakeDownloader{}
	m.dm.ex = fs.FakeExtractor{}
	m.fake = true
	m.fakeError = fakeError
}

func (m *Model) FS() afero.Fs {
	return m.fs
}

func (m *Model) T() *assets.T {
	return m.t
}

func (m *Model) SetT(lang string) {
	t := assets.NewT(lang)
	m.t = &t
}

func (m *Model) OpenWebview(opts WebviewRequest) {
	m.wvCh <- opts
}

func (m *Model) PlayBackgroundMusic(value bool) {
	if value {
		m.player.Pause()
	} else {
		m.player.Play()
	}
}

func (m *Model) Mode() types.Mode {
	if m.mode == types.ModeUnknown {
		return types.ModeHome
	}
	return m.mode
}

func (m *Model) SetMode(mode types.Mode) {
	logger.Info.Printf("navigation: mode changed to %v", mode)
	m.mode = mode
}

func (m *Model) SetListPosition(pos types.ListPosition) {
	m.listPosition = pos
}

func (m *Model) ListPosition() types.ListPosition {
	return m.listPosition
}

func (m *Model) Sidebar() *SidebarModel {
	return &m.sidebar
}

func (m *Model) DownloadManager() *DownloadManagerModel {
	return &m.dm
}
