/*
 *
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/utils/internal/repository"
	"github.com/google/go-github/v56/github"
	"github.com/xanzy/go-gitlab"
)

func NewGitlabClient(token, baseURL string) (client *gitlab.Client, err error) {
	var options []gitlab.ClientOptionFunc

	if baseURL != "" {
		_, err = url.ParseRequestURI(baseURL)
		if err != nil {
			return nil, consts.ErrParamUrl
		}

		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	if consts.ProxyUrl != "" {
		proxyUrl, _ := url.Parse(consts.ProxyUrl)
		httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 5 * time.Second}
		options = append(options, gitlab.WithHTTPClient(httpClient))
	}

	client, err = gitlab.NewClient(token, options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetGitLabTokenInfo(client *gitlab.Client) (owner string, ownerId int64, tokenType int32, expirationTime time.Time, err error) {
	// get token info
	token, _, err := client.PersonalAccessTokens.GetSinglePersonalAccessToken()
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			return "", 0, 0, time.Time{}, consts.ErrTokenInvalid
		}

		return "", 0, 0, time.Time{}, err
	}

	// check token if is valid
	if token.Revoked || !token.Active {
		return "", 0, 0, time.Time{}, consts.ErrTokenInvalid
	}

	// whether if have scope: api
	var hasApi bool
	for _, scope := range token.Scopes {
		if scope == "api" {
			hasApi = true
		}
	}

	if !hasApi {
		return "", 0, 0, time.Time{}, consts.ErrTokenInvalid
	}

	// check token if expired
	expirationTime, err = time.ParseInLocation("2006-01-02", token.ExpiresAt.String(), consts.TimeZone)
	if err != nil {
		return "", 0, 0, time.Time{}, err
	}

	// get user info by token
	user, _, err := client.Users.GetUser(token.UserID, gitlab.GetUsersOptions{})
	if err != nil {
		return "", 0, 0, time.Time{}, err
	}

	// get namespace info by token
	namespaces, _, err := client.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{})
	if err != nil {
		return "", 0, 0, time.Time{}, err
	}

	// check user if is bot
	if user.Bot {
		// if is bot, then token is group access token (https://docs.gitlab.com/ee/user/group/settings/group_access_tokens.html)

		for _, namespace := range namespaces {
			if namespace.Kind == "group" {
				owner = namespace.Name
				ownerId = int64(namespace.ID)
				tokenType = consts.TokenTypeNumOrganization
				break
			}
		}
		if owner == "" {
			// probably token is project token
			return "", 0, 0, time.Time{}, consts.ErrTokenInvalidType
		}
	} else {
		for _, namespace := range namespaces {
			if namespace.Path == user.Username {
				ownerId = int64(namespace.ID)
				break
			}
		}
		owner = user.Username
		tokenType = consts.TokenTypeNumPersonal
	}

	return owner, ownerId, tokenType, expirationTime, nil
}

func NewGithubClient(token string) (client *github.Client, err error) {
	var httpClient *http.Client

	if consts.ProxyUrl != "" {
		proxyUrl, _ := url.Parse(consts.ProxyUrl)
		httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 5 * time.Second}
	}

	client = github.NewClient(httpClient).WithAuthToken(token)

	return client, nil
}

func GetGitHubTokenInfo(client *github.Client) (owner string, ownerId int64, tokenType int32, expirationTime time.Time, err error) {
	user, res, err := client.Users.Get(context.Background(), "")
	if err != nil {
		if githubErr, ok := err.(*github.ErrorResponse); ok {
			if githubErr.Message == "Bad credentials" {
				return "", 0, 0, time.Time{}, consts.ErrTokenInvalid
			}
		}

		return "", 0, 0, time.Time{}, err
	}

	expirationTimeStr := res.Header.Get("github-authentication-token-expiration")
	expirationTime, err = time.ParseInLocation("2006-01-02 15:04:05 MST", expirationTimeStr, consts.TimeZone)
	if err != nil {
		return "", 0, 0, time.Time{}, err
	}

	return *user.Login, *user.ID, consts.TokenTypeNumPersonal, expirationTime, nil
}

func GetRepoFullUrl(repoType int32, repoUrl, ref, filePid string) string {
	switch repoType {
	case consts.RepositoryTypeNumGitLab:
		return repository.GetRepoFullUrlGitLab(repoUrl, ref, filePid)
	case consts.RepositoryTypeNumGithub:
		return repository.GetRepoFullUrlGitHub(repoUrl, ref, filePid)
	default:
		return ""
	}
}

func ParseRepoUrl(url string) (domain, owner, repoName string, err error) {
	r := regexp.MustCompile(repository.RegRepoURL)
	matches := r.FindStringSubmatch(url)
	if len(matches) != 4 {
		return "", "", "", errors.New("repository path format is incorrect; unable to parse the GitHub URL")
	}

	return matches[1], matches[2], matches[3], nil
}

func ParseRepoFileUrl(repoType int32, url string) (filePid, owner, repoName string, err error) {
	switch repoType {
	case consts.RepositoryTypeNumGitLab:
		return repository.ParseRepoFileUrlGitLab(url)

	case consts.RepositoryTypeNumGithub:
		return repository.ParseRepoFileUrlGitHub(url)

	default:
		return "", "", "", errors.New("invalid repo type")
	}
}

func ValidateTokenForRepoGitLab(client *gitlab.Client, owner, repoName string) (bool, error) {
	project, _, err := client.Projects.GetProject(fmt.Sprintf("%s/%s", owner, repoName), &gitlab.GetProjectOptions{})
	if err != nil {
		return false, err
	}

	if project.Permissions.ProjectAccess != nil {
		return project.Permissions.ProjectAccess.AccessLevel >= 40, nil
	}
	if project.Permissions.GroupAccess != nil {
		return project.Permissions.GroupAccess.AccessLevel >= 40, nil
	}

	return false, nil
}

func ValidateTokenForRepoGitHub(client *github.Client, owner, repoName string) (bool, error) {
	repo, _, err := client.Repositories.Get(context.TODO(), owner, repoName)
	if err != nil {
		return false, err
	}

	if has, ok := repo.Permissions["push"]; has && ok {
		if has1, ok1 := repo.Permissions["pull"]; has1 && ok1 {
			return true, nil
		}
	}

	return false, nil
}
