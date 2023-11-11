/*
 * Copyright 2022 CloudWeGo Authors
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
 */

package repository

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strings"
)

type GitLabApi struct {
	client *gitlab.Client
}

func NewGitLabApi(client *gitlab.Client) *GitLabApi {
	return &GitLabApi{
		client: client,
	}
}

const (
	gitlabURLPrefix  = "https://gitlab.com/"
	regGitLabURL     = `([^\/]+)\/([^\/]+)\/-\/blob\/([^\/]+)\/(.+)`
	regGitLabRepoURL = `/gitlab.com/([^/]+)/([^/]+)"`
)

func (a *GitLabApi) ParseIdlUrl(url string) (filePid, owner, repoName string, err error) {
	var tempPath string

	// determine if it is a GitLab prefix
	if strings.HasPrefix(url, gitlabURLPrefix) {
		tempPath = url[len(gitlabURLPrefix):]
		lastQuestionMarkIndex := strings.LastIndex(tempPath, "?")
		if lastQuestionMarkIndex != -1 {
			tempPath = tempPath[:lastQuestionMarkIndex]
		}
	} else {
		return "", "", "", errors.New("idlPath format wrong, do not have prefix: " + gitlabURLPrefix)
	}

	// using regular expressions to match fields
	regex := regexp.MustCompile(regGitLabURL)
	matches := regex.FindStringSubmatch(tempPath)
	if len(matches) != 5 {
		return "", "", "", errors.New("idlPath format wrong, cannot parse gitlab URL")
	}
	owner = matches[1]
	repoName = matches[2]
	filePid = matches[4]

	return filePid, owner, repoName, nil
}

func (a *GitLabApi) ParseRepoUrl(url string) (owner, repoName string, err error) {
	// verification format
	if !strings.HasPrefix(url, gitlabURLPrefix) {
		return "", "", errors.New("IDL path format is incorrect; it does not have the expected prefix: " + gitlabURLPrefix)
	}

	// Extracting information using regular expressions
	r := regexp.MustCompile(regGithubRepoURL)
	matches := r.FindStringSubmatch(url)

	return matches[1], matches[2], nil
}

func (a *GitLabApi) GetFile(owner, repoName, filePath, ref string) (*File, error) {
	// construct the project ID (pid) by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// retrieve the file content from GitLab repository
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}

	// extract the file name from the file path
	name := filePath
	index := strings.LastIndex(filePath, "/")
	if index != -1 {
		name = name[index+1:]
	}

	// decode the base64 encoded file content
	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    name,
		Content: decodedContent,
	}, nil
}

func (a *GitLabApi) PushFilesToRepository(files map[string][]byte, owner, repoName, branch, commitMessage string) error {
	// pushFilesToRepository implementation for GitLab
	for filePath, content := range files {
		contentStr := string(content)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String(commitMessage),
			Content:       &contentStr,
		}

		// create files in the repository
		_, _, err := a.client.RepositoryFiles.CreateFile(fmt.Sprintf("%s/%s", owner, repoName), filePath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *GitLabApi) GetRepositoryArchive(owner, repoName, ref string) ([]byte, error) {
	// generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// specify the desired archive format
	format := "tar"

	// set archive options
	archiveOptions := &gitlab.ArchiveOptions{
		Format: &format, // Choose the archive format
	}

	// request the archive from the GitLab API
	fileData, _, err := a.client.Repositories.Archive(pid, archiveOptions)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (a *GitLabApi) GetLatestCommitHash(owner, repoName, filePath, ref string) (string, error) {
	// generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// request the file content from the GitLab API
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return "", err
	}

	// extract and return the last commit ID
	return fileContent.LastCommitID, nil
}

func (a *GitLabApi) DeleteDirs(owner, repoName string, folderPaths ...string) error {
	// generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// iterate over the folder paths and delete each one
	for _, folderPath := range folderPaths {
		// attempt to delete the specified folder
		_, err := a.client.RepositoryFiles.DeleteFile(pid, folderPath, &gitlab.DeleteFileOptions{
			Branch:        gitlab.String("main"),             // set the branch where the folder is located
			AuthorEmail:   gitlab.String("test@example.com"), // replace with your email
			AuthorName:    gitlab.String("test"),             // replace with your name
			CommitMessage: gitlab.String(fmt.Sprintf("Delete folder %s", folderPath)),
		})

		// check for errors, but ignore if it's a "file not found" error
		if err != nil && !utils.IsFileNotFoundError(err) {
			return err
		}
	}

	return nil
}

func (a *GitLabApi) AutoCreateRepository(owner, repoName string, isPrivate bool) (string, error) {
	// new repository's path in gitlab
	repoPid := owner + "/" + repoName
	// new repository's URL
	newRepoURL := gitlabURLPrefix + owner + "/" + repoName
	_, _, err := a.client.Projects.GetProject(repoPid, nil)
	if err != nil {
		// if the error is caused by the inability to find a repository with the name, create the repository
		if strings.Contains(err.Error(), "404 Project Not Found") {
			var v gitlab.VisibilityValue
			if isPrivate {
				v = gitlab.PrivateVisibility
			} else {
				v = gitlab.PublicVisibility
			}
			_, _, err := a.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Name:        gitlab.String(repoName),
				Visibility:  &v,
				Description: gitlab.String("generate by cwgo"),
			})
			if err != nil {
				return "", err
			}
			return newRepoURL, nil
		}
		return "", err
	}
	return newRepoURL, nil
}

func (a *GitLabApi) GetRepositoryPrivacy(owner, repoName string) (bool, error) {
	repoPid := owner + repoName
	project, _, err := a.client.Projects.GetProject(repoPid, nil)
	if err != nil {
		return false, err
	}
	if project.Visibility == gitlab.PrivateVisibility {
		return true, nil
	} else {
		return false, nil
	}
}
