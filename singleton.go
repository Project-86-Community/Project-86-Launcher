package eightysix

import (
	ESApp "eightysix/internal/app"

	"github.com/quasilyte/gdata/v2"
)

type debugMode struct {
	Logs bool
}

var (
	TheDebugMode debugMode
	GDataM       *gdata.Manager
	app          *ESApp.App
)
