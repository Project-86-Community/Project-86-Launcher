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
	"fmt"
	"p86l/configs"
	"p86l/internal/github"
	"p86l/internal/log"
	"time"

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
	client *github.Client

	data      *CacheData
	expiresAt time.Time
	onRefresh func(*CacheData)
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

// func (c *CacheModel) Load(fs *file.Filesystem) error {
// 	b, err := fs.Load(filepath.Join(configs.AppName, configs.CacheFile))
// 	if err != nil {
// 		return err
// 	}

// 	var d CacheData
// 	if err := json.Unmarshal(b, &d); err != nil {
// 		return err
// 	}

// 	if err := d.Validate(d.Repo); err != nil {
// 		return err
// 	}
// 	if err := d.Validate(d.PreRepo); err != nil {
// 		return err
// 	}

// 	//c.SetVaild(true)
// 	//c.SetFile(d)

// 	return nil
// }

// -- Setters for CacheModel --

func (c *CacheModel) SetLogger(logger *zerolog.Logger) {
	c.logger = logger
}

func (c *CacheModel) SetData(data *CacheData) {
	c.data = data
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

	if c.onRefresh != nil {
		c.onRefresh(data)
	}
}

func (c *CacheModel) Start() {
	if DisableAPI {
		return
	}

	go func() {
		c.refreshData()

		if c.data == nil {
			ctx := context.Background()
			ratelimit, err := c.Client().GetRateLimit(ctx)
			if err != nil {
				c.logger.Warn().Str("CacheModel", "Start").Err(fmt.Errorf("%w: %w", log.ErrCacheRateLimit, err)).Msg(log.ErrorManager.String())
			} else {
				c.data = &CacheData{
					RateLimit2: ratelimit,
					Releases:   nil,
				}
				c.expiresAt = time.Unix(ratelimit.Reset, 0)
			}
		}

		for {
			time.Sleep(time.Until(c.expiresAt))
			c.refreshData()
		}
	}()
}

func (c *CacheModel) ForceRefresh() {
	if c.data != nil && c.data.Releases != nil {
		c.data.Releases = nil
	}
	c.refreshData()
}
