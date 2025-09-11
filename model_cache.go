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
	"p86l/internal/log"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/rs/zerolog"
)

type CacheData struct {
	RateLimit *github.RateLimits        `json:"ratelimit"`
	Repo      *github.RepositoryRelease `json:"repo"`
	PreRepo   *github.RepositoryRelease `json:"prerelease_repo"`
}

func (c *CacheData) Validate(ghRepo *github.RepositoryRelease) error {
	if ghRepo == nil {
		return log.ErrRepoEmpty
	}
	if ghRepo.GetBody() == "" {
		return log.ErrRepoBodyEmpty
	}
	if ghRepo.GetHTMLURL() == "" {
		return log.ErrRepoHTMLURLEmpty
	}
	if len(ghRepo.Assets) < 1 {
		return log.ErrRepoAssetsEmpty
	}

	return nil
}

type CacheModel struct {
	logger    *zerolog.Logger
	data      *CacheData
	expiresAt time.Time
	onRefresh func(*CacheData)
}

// -- Getters for CacheModel --

func (c *CacheModel) Data() *CacheData {
	return c.data
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
	githubClient := github.NewClient(nil)
	ctx1 := context.Background()
	ctx2 := context.Background()

	limit, _, err := githubClient.RateLimit.Get(ctx1)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrCacheRateLimit, err)
	}

	release, _, err := githubClient.Repositories.GetLatestRelease(ctx2, configs.RepoOwner, configs.RepoName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrCacheLatest, err)
	}

	data := &CacheData{
		RateLimit: limit,
		Repo:      release,
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
	c.expiresAt = data.RateLimit.Core.Reset.Time

	if c.onRefresh != nil {
		c.onRefresh(data)
	}
}

func (c *CacheModel) Start() {
	c.refreshData()

	go func() {
		for {
			time.Sleep(time.Until(c.expiresAt))
			c.refreshData()
		}
	}()
}

func (c *CacheModel) ForceRefresh() {
	c.refreshData()
}
