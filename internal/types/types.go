// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package types

/*
ENUM(
unknown = ""
home
settings
about
logs
)
*/
type Mode string

// ENUM(left, right)
type SidebarPosition int

// ENUM(top, bottom)
type ListPosition int

// ENUM(folders, root, versions, logs)
type Folder int

// ENUM(help, cache, report, logs, website, github, discord, patreon, about)
type Helps int

// ENUM(download, extract, done, error)
type Phase int
