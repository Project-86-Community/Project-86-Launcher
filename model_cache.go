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
	"encoding/json"
	"p86l/internal/github"
	"p86l/internal/log"
	"sync"
	"time"

	translator "github.com/Conight/go-googletrans"
)

type CacheFile struct {
	Releases     *github.LatestReleases `json:"releases"`
	RateLimit    *github.RateLimitCore  `json:"rate_limit"`
	LastUpdated  time.Time              `json:"last_updated"`
	ReleasesAge  time.Time              `json:"releases_age"`   // When releases were last fetched
	RateLimitAge time.Time              `json:"rate_limit_age"` // When rate limit was last fetched

	ChangelogTranslation string `json:"-"`
}

type Cache struct {
	mu   sync.RWMutex
	file CacheFile
}

func NewCache(initial CacheFile) *Cache {
	return &Cache{
		file: initial,
	}
}

func (c *Cache) RateLimit() *github.RateLimitCore {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.file.RateLimit
}

func (c *Cache) Releases() *github.LatestReleases {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.file.Releases
}

func (c *Cache) SetRateLimit(rl *github.RateLimitCore) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.file.RateLimit = rl
	c.file.RateLimitAge = time.Now()
	c.file.LastUpdated = time.Now()
}

func (c *Cache) SetReleases(lr *github.LatestReleases) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.file.Releases = lr
	c.file.ReleasesAge = time.Now()
	c.file.LastUpdated = time.Now()
}

// GetReleasesAge returns how old the releases data is
func (c *Cache) ReleasesAge() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.file.ReleasesAge.IsZero() {
		return time.Duration(0)
	}
	return time.Since(c.file.ReleasesAge)
}

// RateLimitAge returns how old the rate limit data is
func (c *Cache) RateLimitAge() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.file.RateLimitAge.IsZero() {
		return time.Duration(0)
	}
	return time.Since(c.file.RateLimitAge)
}

func (c *Cache) ChangelogTranslation() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.file.ChangelogTranslation
}

func (c *Cache) SetChangelogTranslation(body string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.file.ChangelogTranslation = body
}

// -- common --

func (c *Cache) Get() CacheFile {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.file
}

func (c *Cache) Update(fn func(*CacheFile)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(&c.file)
}

func (m *Model) loadCache() error {
	if !m.fs.Exist(m.cachePath) {
		m.logger.Info().Str(log.Lifecycle, "cache file does not exist, using empty").Msg(log.FileManager.String())
		m.cache.Update(func(cf *CacheFile) {
			cf.LastUpdated = time.Now()
		})
		return nil
	}

	jsonData, err := m.fs.Load(m.cachePath)
	if err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to load cache").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	var cache CacheFile
	if err := json.Unmarshal(jsonData, &cache); err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to unmarshal cache").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	m.cache.Update(func(cf *CacheFile) {
		*cf = cache
	})

	m.logger.Info().Str(log.Lifecycle, "cache loaded successfully").Time("last_updated", cache.LastUpdated).Msg(log.FileManager.String())
	return nil
}

func (m *Model) saveCache() error {
	cacheData := m.cache.Get()

	jsonData, err := json.MarshalIndent(cacheData, "", "	")
	if err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to marshal cache").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	if err := m.fs.Save(m.cachePath, jsonData); err != nil {
		m.logger.Warn().Str(log.Lifecycle, "failed to save cache").Err(err).Msg(log.ErrorManager.String())
		return err
	}

	m.logger.Info().Str(log.Lifecycle, "cache saved successfully").Msg(log.FileManager.String())
	return nil
}

func (m *Model) Cache() *Cache {
	return m.cache
}

// -- commands --

type TranslateChangelogCommand struct {
	model            *Model
	origin, dest     string
	googleTranslator *translator.Translator
}

func (t TranslateChangelogCommand) Execute(m *Model) {
	t.model.cache.SetChangelogTranslation("...")

	result, err := t.googleTranslator.Translate(t.origin, "en", t.dest)
	if err != nil {
		m.logger.Warn().Err(err).Msg(log.ErrorManager.String())
		t.model.cache.SetChangelogTranslation(T("errors.translate_fail"))
		return
	}

	t.model.cache.SetChangelogTranslation(result.Text)
	t.model.handleUIRefresh()
}

func (m *Model) Translate(origin, dest string) {
	m.logger.Info().Str(log.Lifecycle, "translating changelog...").Msg(log.NetworkManager.String())
	m.commandChan <- TranslateChangelogCommand{
		model:            m,
		origin:           origin,
		dest:             dest,
		googleTranslator: m.googleTranslator,
	}
}

type ResetCacheCommand struct{}

func (r ResetCacheCommand) Execute(m *Model) {
	m.cache.Update(func(df *CacheFile) {
		*df = CacheFile{}
	})

	m.cacheResetCommandChan <- struct{}{}

	m.handleUIRefresh()
}

func (m *Model) ResetCacheAsync() {
	m.commandChan <- ResetCacheCommand{}
}
