/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86 for managing game files.
 * Copyright (C) 2025 Project 86 Community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package app

import (
	"net"
	"os"
	"p86l"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
	"github.com/rs/zerolog"
)

type modelKey int

const (
	modelKeyModel modelKey = iota
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	sidebar    Sidebar

	model p86l.Model
}

func (r *Root) SetListener(listener net.Listener) {
	r.model.SetListener(listener)
}

func (r *Root) SetLog(logger zerolog.Logger, logFile *os.File) {
	log := r.model.Log()
	log.SetLogger(logger)
	log.SetLogFile(logFile)
}

func (r *Root) Model(key any) any {
	switch key {
	case modelKeyModel:
		return &r.model
	default:
		return nil
	}
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.sidebar)
}

func (r *Root) Build(context *guigui.Context) error {
	context.SetBounds(&r.background, context.Bounds(r), r)

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * u),
			layout.FlexibleSize(1),
		},
	}
	context.SetBounds(&r.sidebar, gl.CellBounds(0, 0), r)

	return nil
}

func (r *Root) Close() error {
	return r.model.Close()
}
