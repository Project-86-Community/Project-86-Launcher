// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package service

import (
	"os/exec"
	"sync"

	"p86l/internal/logger"
)

type launchService struct {
	version    Version
	proc       *exec.Cmd
	mu         sync.Mutex
	runningTag string
	runningOS  string
}

func NewLaunchService() LaunchService {
	return &launchService{}
}

func (s *launchService) CurrentVersion() Version {
	return s.version
}

func (s *launchService) SetCurrentVersion(ver Version) {
	s.version = ver
}

func (s *launchService) IsRunningVersion(tag, os string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.proc != nil && s.proc.ProcessState == nil && s.runningTag == tag && s.runningOS == os
}

func (s *launchService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.proc != nil && s.proc.ProcessState == nil
}

func (s *launchService) Launch(ver Version, rebuildFn func()) error {
	logger.Info.Printf("launching version: %s (%s)", ver.Tag, ver.OS)
	cmd := exec.Command(ver.Executable)

	s.mu.Lock()
	s.proc = cmd
	s.runningTag = ver.Tag
	s.runningOS = ver.OS
	s.mu.Unlock()

	if err := cmd.Start(); err != nil {
		logger.Error.Printf("launch failed [%s]: %v", ver.Tag, err)
		return err
	}
	logger.Info.Printf("launched PID %d [%s %s]", cmd.Process.Pid, ver.Tag, ver.OS)

	go func() {
		_ = cmd.Wait()
		logger.Info.Printf("process exited [%s %s]", ver.Tag, ver.OS)
		s.mu.Lock()
		s.runningTag = ""
		s.runningOS = ""
		s.mu.Unlock()
		if rebuildFn != nil {
			rebuildFn()
		}
	}()

	return nil
}

func (s *launchService) Kill() error {
	s.mu.Lock()
	proc := s.proc
	s.mu.Unlock()

	if proc != nil && proc.Process != nil {
		logger.Info.Printf("killing process PID %d [%s]", proc.Process.Pid, s.version.Tag)
		if err := proc.Process.Kill(); err != nil {
			logger.Warn.Printf("kill failed: %v", err)
			return err
		}
	}
	return nil
}
