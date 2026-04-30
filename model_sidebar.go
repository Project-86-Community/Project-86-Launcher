package p86l

import (
	"os/exec"
	"p86l/internal/logger"
	"p86l/internal/types"
	"sync"

	"github.com/guigui-gui/guigui"
)

type SidebarModel struct {
	version Version

	proc       *exec.Cmd
	procMu     sync.Mutex
	runningTag string // tag of the currently running version
	runningOS  string // OS of the currently running version

	sidebarPosition types.SidebarPosition
}

func (m *SidebarModel) Version() Version {
	return m.version
}

func (m *SidebarModel) SetVersion(ver Version) {
	m.version = ver
}

func (m *SidebarModel) IsRunningVersion(tag, os string) bool {
	m.procMu.Lock()
	defer m.procMu.Unlock()
	return m.proc != nil && m.proc.ProcessState == nil && m.runningTag == tag && m.runningOS == os
}

func (m *SidebarModel) IsRunning() bool {
	m.procMu.Lock()
	defer m.procMu.Unlock()
	return m.proc != nil && m.proc.ProcessState == nil
}

func (m *SidebarModel) Launch(fake bool, widget guigui.Widget) {
	if fake {
		logger.Info.Printf("Launch: fake mode, skipping launch for %s", m.version.Tag)
		return
	}

	logger.Info.Printf("launching version: %s (%s)", m.version.Tag, m.version.OS)
	cmd := exec.Command(m.version.Executable)
	m.procMu.Lock()
	m.proc = cmd
	m.runningTag = m.version.Tag
	m.runningOS = m.version.OS
	m.procMu.Unlock()
	if err := cmd.Start(); err != nil {
		logger.Error.Printf("launch failed [%s]: %v", m.version.Tag, err)
		return
	}
	logger.Info.Printf("launched PID %d [%s %s]", cmd.Process.Pid, m.version.Tag, m.version.OS)
	// Wait in background so ProcessState gets set when it exits.
	go func() {
		_ = cmd.Wait()
		logger.Info.Printf("process exited [%s %s]", m.version.Tag, m.version.OS)
		m.procMu.Lock()
		m.runningTag = ""
		m.runningOS = ""
		m.procMu.Unlock()
		guigui.RequestRebuild(widget)
	}()
}

func (m *SidebarModel) Kill(fake bool) {
	if fake {
		logger.Info.Printf("Kill: fake mode, skipping kill for %s", m.version.Tag)
		return
	}

	m.procMu.Lock()
	proc := m.proc
	m.procMu.Unlock()
	if proc != nil && proc.Process != nil {
		logger.Info.Printf("killing process PID %d [%s]", proc.Process.Pid, m.version.Tag)
		if err := proc.Process.Kill(); err != nil {
			logger.Warn.Printf("kill failed: %v", err)
		}
	}
}

func (m *SidebarModel) SetSidebarPosition(pos types.SidebarPosition) {
	m.sidebarPosition = pos
}

func (m *SidebarModel) SidebarPosition() types.SidebarPosition {
	return m.sidebarPosition
}
