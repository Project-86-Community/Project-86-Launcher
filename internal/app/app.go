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

package app

import (
	"context"
	"net/http"
	"p86l/internal/cache"
	"p86l/internal/data"
	"p86l/internal/debug"
	"p86l/internal/file"
	"time"

	"github.com/google/go-github/v69/github"
)

type App struct {
	isInternet bool
	//Errs       []error

	Debug *debug.Debug
	FS    *file.AppFS
	Data  *data.Data
	Cache *cache.Cache
}

/*
func (a *App) Error(err error) error {
	return errors.New(err.Error())
}

func (a *App) PopupError(err error) {
	newErr := errors.Wrap(err, "Popup error")

	for _, appErr := range a.Errs {
		if appErr.Error() == newErr.Error() {
			return
		}
	}

	log.Error().Stack().Err(newErr).Send()
	a.Errs = append(a.Errs, newErr)
}*/

func (a *App) IsInternet() bool {
	return a.isInternet
}

func (a *App) isInternetReachable() bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 204
}

func (a *App) UpdateInternet() {
	if a.isInternetReachable() {
		a.isInternet = true
	} else {
		a.isInternet = false
	}
}

func (a *App) Update(githubClient *github.Client, context context.Context) {
	// if a.isInternet {
	// 	err := a.Cache.InitChangelog(githubClient, context)
	// 	if err != nil {
	// 		a.PopupError(err)
	// 	}
	// }
}

// func RequestChangelog() (types.Changelog, error) {
// 	changelogData := types.Changelog{}
//
// 	release, _, err := content.GithubClient.Repositories.GetLatestRelease(content.GithubContext, content.RepoOwner, content.RepoName)
// 	if err != nil {
// 		return changelogData, err
// 	}
//
// 	fmt.Println("INTERNET CALL")
//
// 	changelogData.Body = release.GetBody()
// 	changelogData.URL = release.GetHTMLURL()
// 	changelogData.Timestamp = time.Now()
// 	changelogData.ExpiresIn = time.Hour
//
// 	return changelogData, nil
// }

// func WrapText(context *guigui.Context, input string, maxWidth int) string {
// 	charWidthDivisor := 6.2 * context.AppScale()
// 	charCount := int(float64(maxWidth) / charWidthDivisor)
// 	input = strings.ReplaceAll(input, "\r\n", "\n")
//
// 	lines := strings.Split(input, "\n")
// 	var result []string
//
// 	for _, line := range lines {
// 		if len(line) == 0 {
// 			result = append(result, "")
// 			continue
// 		}
//
// 		words := strings.Fields(line)
// 		if len(words) == 0 {
// 			result = append(result, "")
// 			continue
// 		}
//
// 		var currentLine string
//
// 		for _, word := range words {
// 			if len(word) > charCount {
// 				if len(currentLine) > 0 {
// 					result = append(result, currentLine)
// 					currentLine = ""
// 				}
//
// 				for i := 0; i < len(word); i += charCount {
// 					end := i + charCount
// 					if end > len(word) {
// 						end = len(word)
// 					}
// 					result = append(result, word[i:end])
// 				}
// 			} else {
// 				if len(currentLine) == 0 {
// 					currentLine = word
// 				} else if len(currentLine)+1+len(word) <= charCount {
// 					currentLine += " " + word
// 				} else {
// 					result = append(result, currentLine)
// 					currentLine = word
// 				}
// 			}
// 		}
//
// 		if len(currentLine) > 0 {
// 			result = append(result, currentLine)
// 		}
// 	}
//
// 	return strings.Join(result, "\n")
// }
