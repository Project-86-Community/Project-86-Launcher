// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package open

import "os/exec"

func openCmd(input string) *exec.Cmd {
	return exec.Command("open", input)
}
