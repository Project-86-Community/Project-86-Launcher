// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

//go:build !windows && !darwin

package open

import "os/exec"

func openCmd(input string) *exec.Cmd {
	return exec.Command("xdg-open", input)
}
