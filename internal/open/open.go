// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

// Package open provides a cross-platform way to open files, directories, and
// URLs using the OS's default application. It replaces the external
// skratchdot/open-golang dependency with a minimal built-in implementation.
package open

// Start opens a file, directory, or URL using the OS's default application.
// It does not wait for the command to complete.
func Start(input string) error {
	return openCmd(input).Start()
}
