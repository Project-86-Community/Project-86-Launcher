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

package p86l

import (
	"context"
	"encoding/json"
	"fmt"
	"p86l/configs"
	"p86l/internal/file"
	"p86l/internal/github"
	"p86l/internal/log"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
)

type CacheData struct {
	RateLimit2 *github.RateLimitCore  `json:"ratelimit_2"`
	Releases   *github.LatestReleases `json:"releases"`
}

func (c *CacheData) Validate(ghRepo *github.RepositoryRelease) error {
	if ghRepo == nil {
		return log.ErrRepoEmpty
	}
	if ghRepo.Body == "" {
		return log.ErrRepoBodyEmpty
	}
	if len(ghRepo.Assets) < 1 {
		return log.ErrRepoAssetsEmpty
	}

	return nil
}

type CacheModel struct {
	logger *zerolog.Logger
	fs     *file.Filesystem
	client *github.Client

	data      *CacheData
	expiresAt time.Time
}

// -- Getters for CacheModel --

func (c *CacheModel) Client() *github.Client {
	if c.client == nil {
		c.client = github.NewClient(github.Config{})
	}
	return c.client
}

func (c *CacheModel) Data() *CacheData {
	return c.data
}

func (c *CacheModel) ExpiresAt() time.Time {
	return c.expiresAt
}

func (c *CacheModel) Path() string {
	return filepath.Join(configs.AppName, configs.FileCache)
}

func (c *CacheModel) ExpireTimeFormatted() string {
	if c.fs.Exist(c.Path()) {
		if cacheData := c.data; cacheData != nil && cacheData.RateLimit2 != nil {
			return fmt.Sprintf(
				"%d / %d - requests - %s",
				cacheData.RateLimit2.Remaining,
				cacheData.RateLimit2.Limit,
				humanize.RelTime(time.Now(), c.expiresAt, "remaining", "ago"),
			)
		}
	}

	return "..."
}

func (c *CacheModel) GameVersionText(value bool) string {
	if c.data == nil && c.data.Releases == nil {
		return "..."
	}
	if value {
		_, debugGameAsset := GetAssets(c.data.Releases.PreRelease.Assets)
		if debugGameAsset == nil {
			return "..."
		}
		return fmt.Sprintf("%d", debugGameAsset.DownloadCount)
	} else {
		gameAsset, _ := GetAssets(c.data.Releases.Stable.Assets)
		if gameAsset == nil {
			return "..."
		}
		return fmt.Sprintf("%d", gameAsset.DownloadCount)
	}
}

func (c *CacheModel) ChangelogText(value bool) string {
	if c.data == nil && c.data.Releases == nil {
		return "..."
	}
	if value {
		return fmt.Sprintf("%s\n\n%s", c.data.Releases.PreRelease.Name, c.data.Releases.PreRelease.Body)
	} else {
		return fmt.Sprintf("%s\n\n%s", c.data.Releases.Stable.Name, c.data.Releases.Stable.Body)
	}
}

func (c *CacheModel) Load() error {
	b, err := c.fs.Load(c.Path())
	if err != nil {
		return err
	}

	var d *CacheData
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}

	c.data = d

	// if err := d.Validate(d.Repo); err != nil {
	// 	return err
	// }
	// if err := d.Validate(d.PreRepo); err != nil {
	// 	return err
	// }

	//c.SetVaild(true)
	//c.SetFile(d)

	return nil
}

// -- Setters for CacheModel --

func (c *CacheModel) SetLogger(logger *zerolog.Logger) {
	c.logger = logger
}

func (c *CacheModel) SetFS(fs *file.Filesystem) {
	c.fs = fs
}

func (c *CacheModel) Save() error {
	b, err := json.Marshal(c.data)
	if err != nil {
		return err
	}

	err = c.fs.Save(c.Path(), b)
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheModel) Remove() error {
	err := c.fs.Remove(c.Path())
	if err != nil {
		return err
	}

	return nil
}

// -- common --

func (c *CacheModel) fetchData() (*CacheData, error) {
	ctx1 := context.Background()
	ctx2 := context.Background()

	release, err := c.Client().GetLatestReleases(ctx1, configs.RepoOwner, configs.RepoName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrCacheLatest, err)
	}

	ratelimit, err := c.Client().GetRateLimit(ctx2)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrCacheRateLimit, err)
	}

	data := &CacheData{
		RateLimit2: ratelimit,
		Releases:   release,
	}
	return data, nil
}

func (c *CacheModel) refreshData() {
	data, err := c.fetchData()
	if err != nil {
		c.logger.Warn().Str("CacheModel", "refreshData").Err(err).Msg(log.ErrorManager.String())
		return
	}

	c.data = data
	c.expiresAt = time.Unix(data.RateLimit2.Reset, 0)

	err = c.Save()
	if err != nil {
		c.logger.Warn().Str("CacheModel", "refreshData").Err(fmt.Errorf("failed to save cache: %w", err)).Msg(log.ErrorManager.String())
	} else {
		c.logger.Info().Str("CacheModel.refreshData", "saved latest cache").Msg(log.AppManager.String())
	}
}

func (c *CacheModel) Start() {
	// Loads saved data
	if c.fs.Exist(c.Path()) {
		err := c.Load()
		c.logger.Info().Str("CacheModel.Start", "loading cache file").Msg(log.AppManager.String())
		if err != nil {
			c.logger.Warn().Str("CacheModel", "Start").Err(fmt.Errorf("cache corrupted: %w", err)).Msg(log.ErrorManager.String())
		}
		// If cache is expired, set ratelimit to nil.
		if c.data != nil && c.data.RateLimit2 != nil && time.Unix(c.data.RateLimit2.Reset, 0).Before(time.Now()) {
			c.logger.Info().Str("CacheModel.Start", "cache is expired").Msg(log.AppManager.String())
			c.data.RateLimit2 = nil
			c.refreshData()
		}
	}

	if DisableAPI {
		c.logger.Info().Str("CacheModel", "API is disabled now").Msg(log.AppManager.String())
		return
	}

	// Refresh data when, on startup and there is no cache.
	if !c.fs.Exist(c.Path()) {
		c.logger.Info().Str("CacheModel.Start", "no cache file found, getting new cache").Msg(log.AppManager.String())
		c.refreshData()
	}

	// Gets the ratelimit when api is ratelimited.
	if c.data == nil {
		ctx := context.Background()
		ratelimit, err := c.Client().GetRateLimit(ctx)
		if err != nil {
			c.logger.Warn().Str("CacheModel", "Start").Err(fmt.Errorf("%w: %w", log.ErrCacheRateLimit, err)).Msg(log.ErrorManager.String())
		} else {
			c.logger.Info().Str("CacheModel.Start", "API is ratelimited").Msg(log.AppManager.String())
			c.data = &CacheData{
				RateLimit2: ratelimit,
				Releases:   nil,
			}
			c.expiresAt = time.Unix(ratelimit.Reset, 0)
		}
	}

	// Refresh data when sleep is over.
	for {
		// Wait for user to do ForceRefresh,
		// this situation only happens when,
		// 1. timestamp is expired.
		// 2. there is no cache and internet to fetch data.
		if c.data.RateLimit2 == nil {
			continue
		}
		time.Sleep(time.Until(c.expiresAt))
		c.refreshData()
	}
}

// Removes data, but keeps ratelimit saved locally in app.
func (c *CacheModel) ForceRefresh() {
	if c.data != nil && c.data.Releases != nil {
		c.data.Releases = nil
	}
	if err := c.Remove(); err != nil {
		c.logger.Warn().Str("CacheModel", "ForceRefresh").Err(fmt.Errorf("failed to remove cache: %w", err))
	}
	c.logger.Info().Str("CacheModel.ForceRefresh", "force refreshing cache").Msg(log.AppManager.String())
	c.refreshData()
}
