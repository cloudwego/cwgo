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
	if strings.HasPrefix(url, gitlabURLPrefix) {
		tempPath = url[len(gitlabURLPrefix):]
		lastQuestionMarkIndex := strings.LastIndex(tempPath, "?")
		if lastQuestionMarkIndex != -1 {
			tempPath = tempPath[:lastQuestionMarkIndex]
		}
	} else {
		return "", "", "", errors.New("idlPath format wrong,do not have prefix: " + gitlabURLPrefix)
	}
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
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}

	name := filePath
	index := strings.LastIndex(filePath, "/")
	if index != -1 {
		name = name[index+1:]
	}

	// Decoding the content of base64 encoded files
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
	// Implement PushFilesToRepository for GitLab
	for filePath, content := range files {
		contentStr := string(content)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String(commitMessage),
			Content:       &contentStr,
		}

		_, _, err := a.client.RepositoryFiles.CreateFile(fmt.Sprintf("%s/%s", owner, repoName), filePath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *GitLabApi) GetRepositoryArchive(owner, repoName, format, ref string) ([]byte, error) {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	archiveOptions := &gitlab.ArchiveOptions{
		Format: &format, // Choose the archive format
	}

	fileData, _, err := a.client.Repositories.Archive(pid, archiveOptions)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (a *GitLabApi) GetLatestCommitHash(owner, repoName, filePath, ref string) (string, error) {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	fileContent, _, err := a.client.RepositoryFiles.GetFile(pid, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return "", err
	}
	return fileContent.LastCommitID, nil
}

func (a *GitLabApi) DeleteDirs(owner, repoName string, folderPaths ...string) error {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	for _, folderPath := range folderPaths {
		_, err := a.client.RepositoryFiles.DeleteFile(pid, folderPath, &gitlab.DeleteFileOptions{
			Branch:        gitlab.String("main"),             // 设置要删除的文件夹所在的分支
			AuthorEmail:   gitlab.String("test@example.com"), // 替换为您的邮箱
			AuthorName:    gitlab.String("test"),             // 替换为您的名称
			CommitMessage: gitlab.String(fmt.Sprintf("Delete folder %s", folderPath)),
		})
		if err != nil && !utils.IsFileNotFoundError(err) {
			return err
		}
	}

	return nil
}
