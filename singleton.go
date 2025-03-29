package p86l

import (
	ESApp "p86l/internal/app"

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
