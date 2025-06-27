
package file

import (
	"bytes"
	"encoding/gob"
	"errors"
	"p86l/internal/debug"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog/log"
)

type Data struct {
	Locale      string
	AppScale    int
	ColorMode   guigui.ColorMode
	GameVersion string
}

func (d *Data) Log() {
	log.Info().Any("Translation", d.Locale).Msg("FileManager")
	log.Info().Any("Scaling", d.AppScale).Msg("FileManager")
	log.Info().Any("Theme", d.ColorMode).Msg("FileManager")
	if d.GameVersion == "" {
		return
	}
	log.Info().Any("Game Version", d.GameVersion).Msg("FileManager")
}

type Cache struct {
	Repo      *github.RepositoryRelease
	Timestamp time.Time
	ExpiresIn time.Duration
}

func (c *Cache) Log() {
	log.Info().Any("Changelog", c.Repo.GetBody()).Any("Timestamp", c.Timestamp).Any("ExpiresIn", c.ExpiresIn).Msg("FileManager")
}

func (c *Cache) Validate(appDebug *debug.Debug) *debug.Error {
	if c.Repo == nil {
		return appDebug.New(errors.New("repo is empty"), debug.CacheError, debug.ErrCacheInvalid)
	}

	if c.Repo.GetBody() == "" {
		return appDebug.New(errors.New("body is empty"), debug.CacheError, debug.ErrCacheBodyInvalid)
	}

	if c.Repo.GetHTMLURL() == "" {
		return appDebug.New(errors.New("URL is empty"), debug.CacheError, debug.ErrCacheURLInvalid)
	}
	if len(c.Repo.Assets) < 1 {
		return appDebug.New(errors.New("assets are empty"), debug.CacheError, debug.ErrCacheAssetsInvalid)
	}

	return nil
}

func (a *AppFS) DecodeData(appDebug *debug.Debug, b []byte) (Data, *debug.Error) {
	var d Data
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&d); err != nil {
		return d, appDebug.New(err, debug.FSError, debug.ErrDataLoad)
	}
	return d, nil
}

func (a *AppFS) DecodeCache(appDebug *debug.Debug, b []byte) (Cache, *debug.Error) {
	var c Cache
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&c); err != nil {
		return c, appDebug.New(err, debug.FSError, debug.ErrCacheLoad)
	}
	return c, nil
}
