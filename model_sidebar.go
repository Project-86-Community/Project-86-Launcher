package p86l

import "p86l/internal/types"

type SidebarModel struct {
	sidebarPosition types.SidebarPosition
}

func (m *SidebarModel) SetSidebarPosition(pos types.SidebarPosition) {
	m.sidebarPosition = pos
}

func (m *SidebarModel) SidebarPosition() types.SidebarPosition {
	return m.sidebarPosition
}
