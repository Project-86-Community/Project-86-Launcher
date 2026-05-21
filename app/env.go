// Copyright (c) 2026 Project 86 Community
// SPDX-License-Identifier: GPL-3.0-only

package app

import (
	"fmt"

	"github.com/guigui-gui/guigui"
)

// envMust retrieves an environment value from the widget tree and type-asserts it to T.
// It returns the typed value and true if the key was found, or the zero value of T and false otherwise.
func envMust[T any](context *guigui.Context, widget guigui.Widget, key guigui.EnvKey) (T, bool) {
	v, ok := context.Env(widget, key)
	if !ok {
		var zero T
		return zero, false
	}
	t, ok := v.(T)
	if !ok {
		var zero T
		panic(fmt.Sprintf("envMust: type assertion failed for key %v: expected %T, got %T", key, zero, v))
	}
	return t, true
}
