/*
 * SPDX-License-Identifier: GPL-3.0-only
 * SPDX-FileCopyrightText: 2025 Project 86 Community
 *
 * Project-86-Launcher: A Launcher developed for Project-86-Community-Game for managing game files.
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

package p86l

import (
	"context"
	"p86l/configs"
	"p86l/internal/file"
	"p86l/internal/github"
	"p86l/internal/log"
	"path/filepath"
	"sync"
	"time"

	translator "github.com/Conight/go-googletrans"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/rs/zerolog"
)

type Command interface {
	Execute(*Model)
}

type SubModel interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type Model struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	subModels []SubModel

	logger           *zerolog.Logger
	fs               *file.Filesystem
	bgmPlayer        *audio.Player
	googleTranslator *translator.Translator

	uiRefreshFn func()
	syncDataFn  func(m *Model, value bool) error

	isAutoUseDarkmode bool
	isNew             bool
	dataPath          string
	data              *Data

	cachePath string
	cache     *Cache

	commandChan           chan Command
	cacheResetCommandChan chan struct{}

	isAvailStable, isAvailPreRelease bool
	fileAvailability                 map[string]bool
	fileAvailMutex, uiRefreshFnMutex sync.RWMutex
}

func NewModel(logger *zerolog.Logger, fs *file.Filesystem, bgmPlayer *audio.Player) *Model {
	ctx, cancel := context.WithCancel(context.Background())
	return &Model{
		ctx:                   ctx,
		cancel:                cancel,
		logger:                logger,
		fs:                    fs,
		bgmPlayer:             bgmPlayer,
		googleTranslator:      translator.New(),
		subModels:             make([]SubModel, 0),
		dataPath:              filepath.Join(configs.AppName, configs.FileData),
		data:                  NewData(DataFile{}),
		cachePath:             filepath.Join(configs.AppName, configs.FileCache),
		cache:                 NewCache(CacheFile{}),
		commandChan:           make(chan Command, 10),
		cacheResetCommandChan: make(chan struct{}, 1),
		fileAvailability:      make(map[string]bool),
	}
}

func (m *Model) BGMPlayer() *audio.Player {
	return m.bgmPlayer
}

func (m *Model) IsAutoUseDarkmode() bool {
	return m.isAutoUseDarkmode
}

// - SetUIRefreshFn & SetSyncDataFn both are involved in ui refresh.

// SetUIRefreshFn used for resetting ui, when new info needs to be shown.
func (m *Model) SetUIRefreshFn(fn func()) {
	m.uiRefreshFnMutex.Lock()
	defer m.uiRefreshFnMutex.Unlock()
	m.uiRefreshFn = fn
}

// SetSyncDataFn used for syncing data from Model to UI.
func (m *Model) SetSyncDataFn(fn func(m *Model, value bool) error) {
	m.uiRefreshFnMutex.Lock()
	defer m.uiRefreshFnMutex.Unlock()
	m.syncDataFn = fn
}

// SetIsAutoUseDarkmode whether the user's device uses darkmode by default,
// used for ResetDataAsync, to change theme back to automode.
func (m *Model) SetIsAutoUseDarkmode(value bool) {
	m.isAutoUseDarkmode = value
}

// -- common --

func (m *Model) AddSubModel(sub SubModel) {
	m.subModels = append(m.subModels, sub)
}

func (m *Model) Start() {
	logger := m.logger.With().Str(log.UnknownModel.String(), log.MainModel.String()).Logger()

	logger.Info().Str(log.Lifecycle, log.Starting).Msg(log.AppManager.String())

	for _, sub := range m.subModels {
		sub.Start(m.ctx, &m.wg)
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		logger.Info().Str(log.BackgroundLoop, log.Starting).Msg(log.AppManager.String())
		for {
			select {
			case <-m.ctx.Done():
				logger.Info().Str(log.BackgroundLoop, log.Stopped).Msg(log.AppManager.String())

				if err := m.saveData(); err != nil {
					logger.Warn().Str(log.BackgroundLoop, "failed to save data on shutdown").Err(err).Msg(log.ErrorManager.String())
				}

				return
			case cmd := <-m.commandChan:
				cmd.Execute(m)
			}
		}
	}()
}

func (m *Model) OpenPath(path string) {
	m.logger.Info().Str("open path", path).Msg(log.AppManager.String())
	m.fs.Open(filepath.Join(m.fs.Path(), path))
}

func (m *Model) OpenURL(url string) {
	m.logger.Info().Str("open url", url).Msg(log.AppManager.String())
	m.fs.Open(url)
}

// CheckFilesCached returns cached file availability (instant, non-blocking)
func (m *Model) CheckFilesCached(filePath string) (bool, bool) {
	m.fileAvailMutex.RLock()
	defer m.fileAvailMutex.RUnlock()
	available, exists := m.fileAvailability[filePath]
	return available, exists // (isAvailable, wasCached)
}

func (m *Model) handleUIRefresh() {
	m.uiRefreshFnMutex.RLock()
	refreshFn := m.uiRefreshFn
	m.uiRefreshFnMutex.RUnlock()

	if refreshFn != nil {
		refreshFn()
	}
}

func (m *Model) Stop() {
	m.cancel()
	m.wg.Wait()
}

// -- subModels --

type DataSubModel struct {
	model *Model
}

func NewDataSubModel(model *Model) *DataSubModel {
	return &DataSubModel{
		model: model,
	}
}

func (d *DataSubModel) Start(ctx context.Context, wg *sync.WaitGroup) {
	logger := d.model.logger.With().Str(log.UnknownModel.String(), log.DataModel.String()).Logger()

	logger.Info().Str(log.Lifecycle, log.Starting).Msg(log.AppManager.String())

	if err := d.model.loadData(); err != nil {
		logger.Warn().Str(log.Lifecycle, "could not load initial data, using defaults").Err(err).Msg(log.ErrorManager.String())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info().Str(log.BackgroundLoop, log.Starting).Msg(log.AppManager.String())

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info().Str(log.BackgroundLoop, log.Stopped).Msg(log.AppManager.String())
				return
			case <-ticker.C:
				d.checkFiles()
			}
		}
	}()
}

func (d *DataSubModel) updateFilesCache(filePaths ...string) {
	results := make(map[string]bool)
	for _, path := range filePaths {
		results[path] = d.model.fs.Exist(path)
	}

	d.model.fileAvailMutex.Lock()
	for path, available := range results {
		d.model.fileAvailability[path] = available
	}
	d.model.fileAvailMutex.Unlock()
}

func (d *DataSubModel) checkFiles() {
	filesToCheck := []string{
		PathGameStable,
		PathGamePreRelease,
	}

	d.updateFilesCache(filesToCheck...)

	value1, _ := d.model.CheckFilesCached(PathGameStable)
	value2, _ := d.model.CheckFilesCached(PathGamePreRelease)

	if value1 != d.model.isAvailStable || value2 != d.model.isAvailPreRelease {
		d.model.logger.Info().
			Str(log.Lifecycle, "game files availability changed").
			Bool("stable_was", d.model.isAvailStable).
			Bool("stable_now", value1).
			Bool("prerelease_was", d.model.isAvailPreRelease).
			Bool("prerelease_now", value2).
			Msg(log.AppManager.String())

		d.model.isAvailStable = value1
		d.model.isAvailPreRelease = value2

		d.model.handleUIRefresh()
	}
}

const (
	// - Refresh intervals
	defaultRefreshInterval   = time.Minute
	minRefreshInterval       = time.Second * 5
	rateLimitRefreshInterval = 5 * time.Minute
	releasesRefreshInterval  = 30 * time.Minute
)

type CacheSubModel struct {
	logger zerolog.Logger
	model  *Model

	client *github.Client
}

func NewCacheSubModel(model *Model) *CacheSubModel {
	return &CacheSubModel{
		logger: model.logger.With().Str(log.UnknownModel.String(), log.CacheModel.String()).Logger(),
		model:  model,
		client: github.NewClient(),
	}
}

func (c *CacheSubModel) getRefreshInterval() time.Duration {
	resetTime := c.model.Cache().Get().RateLimit.Reset
	if resetTime <= 0 {
		return defaultRefreshInterval
	}

	interval := time.Until(time.Unix(resetTime, 0))
	if interval <= 0 {
		return minRefreshInterval
	}

	return interval
}

func (c *CacheSubModel) Start(ctx context.Context, wg *sync.WaitGroup) {
	c.logger.Info().Str(log.Lifecycle, log.Starting).Msg(log.AppManager.String())

	if err := c.model.loadCache(); err != nil {
		c.logger.Warn().Str(log.Lifecycle, "could not load cache").Err(err).Msg(log.ErrorManager.String())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.logger.Info().Str(log.BackgroundLoop, log.Starting).Msg(log.AppManager.String())

		c.initialFetch(ctx)

		refT := c.getRefreshInterval()
		c.logger.Info().Str("refresh ticker", refT.String()).Msg(log.AppManager.String())

		refreshTicker := time.NewTicker(refT)
		defer refreshTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info().Str(log.BackgroundLoop, log.Stopped).Msg(log.AppManager.String())

				// Save cache before stopping
				if err := c.model.saveCache(); err != nil {
					c.logger.Warn().Str(log.BackgroundLoop, "failed to save cache on shutdown").Err(err).Msg(log.ErrorManager.String())
				}

				return
			case <-refreshTicker.C:
				c.logger.Debug().Msg("time up")
				c.fetchReleases(ctx)
				c.fetchRateLimit(ctx)
				refreshTicker.Reset(c.getRefreshInterval())
			case <-c.model.cacheResetCommandChan:
				c.fetchReleases(ctx)
				c.fetchRateLimit(ctx)
				refreshTicker.Reset(c.getRefreshInterval())
			}
		}
	}()
}

func (c *CacheSubModel) initialFetch(ctx context.Context) {
	cache := c.model.cache
	cacheFile := cache.Get()

	hasRateLimit := cacheFile.RateLimit != nil
	hasReleases := cacheFile.Releases != nil

	rateLimitAge := cache.RateLimitAge()
	releasesAge := cache.ReleasesAge()

	c.logger.Info().
		Str(log.InitialFetch, "cache status").
		Bool("has_releases", hasReleases).
		Bool("has_rate_limit", hasRateLimit).
		Dur("releases_age", releasesAge).
		Dur("rate_limit_age", rateLimitAge).
		Msg(log.AppManager.String())

	// Always fetch rate limit on start (no rate limit on this endpoint)
	if !hasRateLimit || rateLimitAge > rateLimitRefreshInterval {
		c.logger.Info().Str(log.InitialFetch, "fetching initial rate limit").Msg(log.AppManager.String())
		c.fetchRateLimit(ctx)
	}

	// Only fetch releases if we don't have them or they're very old
	if !hasReleases || releasesAge > releasesRefreshInterval {
		c.logger.Info().Str(log.InitialFetch, "fetching initial releases").Msg(log.AppManager.String())
		c.fetchReleases(ctx)
		c.fetchRateLimit(ctx)
	} else {
		c.logger.Info().
			Str(log.InitialFetch, "using cached releases").
			Dur("age", releasesAge).
			Msg(log.AppManager.String())
	}
}

func (c *CacheSubModel) fetchRateLimit(ctx context.Context) {
	c.logger.Info().Str(log.FetchRateLimit, "fetching rate limit").Msg(log.AppManager.String())
	cache := c.model.Cache()

	rl, err := c.client.GetRateLimit(ctx)
	if err != nil {
		c.logger.Warn().Str(log.FetchRateLimit, "failed to fetch rate limit").Err(err).Msg(log.ErrorManager.String())
		return
	}

	c.logger.Info().
		Str(log.FetchRateLimit, "rate limit updated").
		Int("remaining", rl.Remaining).
		Int("limit", rl.Limit).
		Msg(log.AppManager.String())

	cache.SetRateLimit(rl)
	c.model.handleUIRefresh()
}

func (c *CacheSubModel) fetchReleases(ctx context.Context) {
	cache := c.model.Cache()

	ratelimit := cache.Get().RateLimit
	if ratelimit != nil && ratelimit.Remaining < 5 {
		c.logger.Warn().
			Str(log.FetchReleases, "rate limit low, skipping fetch").
			Int("remaining", ratelimit.Remaining).
			Msg(log.ErrorManager.String())
		return
	}

	c.logger.Info().Str(log.FetchReleases, "fetching releases").Msg(log.AppManager.String())

	lr, err := c.client.GetLatestReleases(ctx, configs.RepoOwner, configs.RepoName)
	if err != nil {
		c.logger.Warn().
			Str(log.FetchReleases, "failed to fetch releases").
			Err(err).
			Msg(log.ErrorManager.String())
		return
	}

	c.logger.Info().
		Str(log.FetchReleases, "releases updated").
		Msg(log.AppManager.String())

	cache.SetReleases(lr)
	c.model.handleUIRefresh()
}
