// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package open

import (
	"os"
	"os/exec"
	"path/filepath"
)

var runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")

func openCmd(input string) *exec.Cmd {
	return exec.Command(runDll32, "url.dll,FileProtocolHandler", input)
}
