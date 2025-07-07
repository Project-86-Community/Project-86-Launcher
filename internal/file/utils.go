package file

import (
	"encoding/json"
	"errors"
	pd "p86l/internal/debug"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/hajimehoshi/guigui"
	"github.com/rs/zerolog/log"
)

type Data struct {
	V             int              `json:"v"`
	Locale        string           `json:"locale"`
	AppScale      int              `json:"app_scale"`
	ColorMode     guigui.ColorMode `json:"color_mode"`
	UsePreRelease bool             `json:"use_pre_release"`
	GameVersion   string           `json:"game_version"`
}

func (d *Data) Log() {
	log.Info().Any("Translation", d.Locale).Msg("FileManager")
	log.Info().Any("Scaling", d.AppScale).Msg("FileManager")
	log.Info().Any("Theme", d.ColorMode).Msg("FileManager")
	if d.GameVersion == "" {
		return
	}
	log.Info().Any("Use Pre-release", d.UsePreRelease).Msg("FileManager")
	log.Info().Any("Game Version", d.GameVersion).Msg("FileManager")
}

type Cache struct {
	V         int                       `json:"v"`
	Repo      *github.RepositoryRelease `json:"repo"`
	Timestamp time.Time                 `json:"time_stamp"`
	ExpiresIn time.Duration             `json:"expires_in"`
}

func (c *Cache) Log() {
	log.Info().Any("Changelog", c.Repo.GetBody()).Any("Timestamp", c.Timestamp).Any("ExpiresIn", c.ExpiresIn).Msg("FileManager")
}

func (c *Cache) Validate(appDebug *pd.Debug) *pd.Error {
	if c.Repo == nil {
		return appDebug.New(errors.New("repo is empty"), pd.CacheError, pd.ErrCacheInvalid)
	}

	if c.Repo.GetBody() == "" {
		return appDebug.New(errors.New("body is empty"), pd.CacheError, pd.ErrCacheBodyInvalid)
	}

	if c.Repo.GetHTMLURL() == "" {
		return appDebug.New(errors.New("URL is empty"), pd.CacheError, pd.ErrCacheURLInvalid)
	}
	if len(c.Repo.Assets) < 1 {
		return appDebug.New(errors.New("assets are empty"), pd.CacheError, pd.ErrCacheAssetsInvalid)
	}

	return nil
}

// -- Used to convert data to bytes.

func (a *AppFS) EncodeData(appDebug *pd.Debug, d Data) ([]byte, *pd.Error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, appDebug.New(err, pd.FSError, pd.ErrDataSave)
	}
	return b, nil
}

func (a *AppFS) EncodeCache(appDebug *pd.Debug, c Cache) ([]byte, *pd.Error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, appDebug.New(err, pd.FSError, pd.ErrDataSave)
	}
	return b, nil
}

// -- Used to get data and use it for app.

// Get data.
func (a *AppFS) DecodeData(appDebug *pd.Debug, b []byte) (Data, *pd.Error) {
	var d Data
	err := json.Unmarshal(b, &d)
	if err != nil {
		return d, appDebug.New(err, pd.FSError, pd.ErrDataLoad)
	}
	return d, nil
}

// Get cache.
func (a *AppFS) DecodeCache(appDebug *pd.Debug, b []byte) (Cache, *pd.Error) {
	var c Cache
	err := json.Unmarshal(b, &c)
	if err != nil {
		return c, appDebug.New(err, pd.FSError, pd.ErrCacheLoad)
	}
	return c, nil
}
