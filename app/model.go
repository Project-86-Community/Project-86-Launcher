// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package app

import (
	"p86l/assets"
	"p86l/internal/types"
)

// Model holds pure application UI state
// It contains no service references; services remain as separate env keys.
type Model struct {
	mode         types.Mode
	listPosition types.ListPosition
	sidebarPos   types.SidebarPosition
	t            *assets.T
}

func NewModel() *Model {
	return &Model{
		mode:         types.ModeHome,
		listPosition: types.ListPositionTop,
		sidebarPos:   types.SidebarPositionRight,
	}
}

func (m *Model) Mode() types.Mode {
	return m.mode
}

func (m *Model) SetMode(mode types.Mode) {
	m.mode = mode
}

func (m *Model) ListPosition() types.ListPosition {
	return m.listPosition
}

func (m *Model) SetListPosition(pos types.ListPosition) {
	m.listPosition = pos
}

func (m *Model) SidebarPosition() types.SidebarPosition {
	return m.sidebarPos
}

func (m *Model) SetSidebarPosition(pos types.SidebarPosition) {
	m.sidebarPos = pos
}

func (m *Model) T() *assets.T {
	return m.t
}

func (m *Model) SetT(t *assets.T) {
	m.t = t
}
