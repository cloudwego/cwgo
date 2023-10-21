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
	gitlabURLPrefix = "https://gitlab.com/"
)

func (a *GitLabApi) ParseUrl(url string) (filePid, owner, repoName string, err error) {
	var tempPath string

	// Determine if it is a GitLab prefix
	if strings.HasPrefix(url, gitlabURLPrefix) {
		tempPath = url[len(gitlabURLPrefix):]
		lastQuestionMarkIndex := strings.LastIndex(tempPath, "?")
		if lastQuestionMarkIndex != -1 {
			tempPath = tempPath[:lastQuestionMarkIndex]
		}
	} else {
		return "", "", "", errors.New("idlPath format wrong,do not have prefix: " + gitlabURLPrefix)
	}

	// Using regular expressions to match fields
	regex := regexp.MustCompile(`([^\/]+)\/([^\/]+)\/-\/blob\/([^\/]+)\/(.+)`)
	matches := regex.FindStringSubmatch(tempPath)
	if len(matches) != 5 {
		return "", "", "", errors.New("idlPath format wrong,cannot parse gitlab URL")
	}
	owner = matches[1]
	repoName = matches[2]
	filePid = matches[4]

	return filePid, owner, repoName, nil
}

func (a *GitLabApi) GetFile(owner, repoName, filePath, ref string) (*File, error) {
	// Construct the project ID (pid) by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// Retrieve the file content from GitLab repository
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}

	// Extract the file name from the file path
	name := filePath
	index := strings.LastIndex(filePath, "/")
	if index != -1 {
		name = name[index+1:]
	}

	// Decode the base64 encoded file content
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
	// PushFilesToRepository implementation for GitLab
	for filePath, content := range files {
		contentStr := string(content)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String(commitMessage),
			Content:       &contentStr,
		}

		// Create files in the repository
		_, _, err := a.client.RepositoryFiles.CreateFile(fmt.Sprintf("%s/%s", owner, repoName), filePath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *GitLabApi) GetRepositoryArchive(owner, repoName, ref string) ([]byte, error) {
	// Generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// Specify the desired archive format
	format := "tar"

	// Set archive options
	archiveOptions := &gitlab.ArchiveOptions{
		Format: &format, // Choose the archive format
	}

	// Request the archive from the GitLab API
	fileData, _, err := a.client.Repositories.Archive(pid, archiveOptions)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (a *GitLabApi) GetLatestCommitHash(owner, repoName, filePath, ref string) (string, error) {
	// Generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// Request the file content from the GitLab API
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return "", err
	}

	// Extract and return the last commit ID
	return fileContent.LastCommitID, nil
}

func (a *GitLabApi) DeleteDirs(owner, repoName string, folderPaths ...string) error {
	// Generate the project ID by combining owner and repoName
	pid := fmt.Sprintf("%s/%s", owner, repoName)

	// Iterate over the folder paths and delete each one
	for _, folderPath := range folderPaths {
		// Attempt to delete the specified folder
		_, err := a.client.RepositoryFiles.DeleteFile(pid, folderPath, &gitlab.DeleteFileOptions{
			Branch:        gitlab.String("main"),             // Set the branch where the folder is located
			AuthorEmail:   gitlab.String("test@example.com"), // Replace with your email
			AuthorName:    gitlab.String("test"),             // Replace with your name
			CommitMessage: gitlab.String(fmt.Sprintf("Delete folder %s", folderPath)),
		})

		// Check for errors, but ignore if it's a "file not found" error
		if err != nil && !utils.IsFileNotFoundError(err) {
			return err
		}
	}

	return nil
}
