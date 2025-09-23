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

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"p86l/internal/log"
	"sort"
	"time"
)

type RateLimitCore struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

type RateLimits struct {
	Resources struct {
		Core RateLimitCore `json:"core"`
	} `json:"resources"`
}

type RepositoryRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Prerelease  bool           `json:"prerelease"`
	Body        string         `json:"body"`
	Assets      []ReleaseAsset `json:"assets"`
	CreatedAt   time.Time      `json:"created_at"`
	PublishedAt time.Time      `json:"published_at"`
}

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	DownloadCount      int64  `json:"download_count"`
}

type LatestReleases struct {
	PreRelease *RepositoryRelease `json:"pre_release"`
	Stable     *RepositoryRelease `json:"stable"`
}

type Config struct {
	Timeout time.Duration
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 0 * time.Second,
		},
	}
}

func (c *Client) BaseURL() string {
	return "https://api.github.com"
}

func (c *Client) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	url := c.BaseURL() + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrGithubRequestNew, err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrGithubRequestDo, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d, %s", log.ErrGithubRequestStatus, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", log.ErrGithubRequestBodyRead, err)
	}

	return body, nil
}

func (c *Client) GetRateLimit(ctx context.Context) (*RateLimitCore, error) {
	data, err := c.doRequest(ctx, "GET", "/rate_limit")
	if err != nil {
		return nil, err
	}

	var rateLimitResp RateLimits
	if err := json.Unmarshal(data, &rateLimitResp); err != nil {
		return nil, fmt.Errorf("parsing rate limit data: %w", err)
	}

	return &rateLimitResp.Resources.Core, nil
}

func (c *Client) getReleases(ctx context.Context, owner, repo string) ([]RepositoryRelease, error) {
	path := fmt.Sprintf("/repos/%s/%s/releases", owner, repo)
	data, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var releases []RepositoryRelease
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil, fmt.Errorf("parsing releases data: %w", err)
	}

	// Sort by published date (newest first)
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedAt.After(releases[j].PublishedAt)
	})

	return releases, nil
}

func (c *Client) GetLatestReleases(ctx context.Context, owner, repo string) (*LatestReleases, error) {
	releases, err := c.getReleases(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	result := &LatestReleases{}

	for i := range releases {
		release := &releases[i]

		if release.Prerelease && result.PreRelease == nil {
			result.PreRelease = release
		} else if !release.Prerelease && result.Stable == nil {
			result.Stable = release
		}

		// Exit early if we found both
		if result.PreRelease != nil && result.Stable != nil {
			break
		}
	}

	return result, nil
}
