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

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"p86l/configs"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/quasilyte/gdata/v2"
	"github.com/rs/zerolog/log"
)

type Changelog struct {
	Body      string
	URL       string
	Timestamp time.Time
	ExpiresIn time.Duration
}

type Cache struct {
	Changelog *Changelog

	GDataM *gdata.Manager
}

func (c *Cache) saveChangelog() error {
	if c.Changelog == nil {
		return errors.New("Changelog not found")
	} else {
		changelogBytes, err := json.Marshal(c.Changelog)
		if err != nil {
			return err
		}
		if err := c.GDataM.SaveObjectProp(configs.Cache, configs.ChangelogFile, changelogBytes); err != nil {
			return err
		}
		return nil
	}
}

func (c *Cache) InitChangelog(githubClient *github.Client, context context.Context) error {
	if c.GDataM.ObjectPropExists(configs.Cache, configs.ChangelogFile) {
		changelogJSON, err := c.GDataM.LoadObjectProp(configs.Cache, configs.ChangelogFile)
		if err != nil {
			return err
		}
		changelogData := &Changelog{}
		err = json.Unmarshal(changelogJSON, &changelogData)
		if err != nil {
			return err
		}
		c.Changelog = changelogData
	} else {
		changelogData, err := c.RequestChangelog(githubClient, context)
		if err != nil {
			return err
		}
		c.Changelog = &changelogData

		err = c.saveChangelog()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) RequestChangelog(githubClient *github.Client, context context.Context) (Changelog, error) {
	changelogData := Changelog{}

	release, _, err := githubClient.Repositories.GetLatestRelease(context, configs.RepoOwner, configs.RepoName)
	if err != nil {
		return changelogData, err
	}

	log.Info().Msg("INTERNET CALL")

	changelogData.Body = release.GetBody()
	changelogData.URL = release.GetHTMLURL()
	changelogData.Timestamp = time.Now()
	changelogData.ExpiresIn = time.Hour

	return changelogData, nil
}
