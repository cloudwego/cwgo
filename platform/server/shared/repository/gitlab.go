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
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type GitLabApi struct {
	Client *gitlab.Client
}

func (gl *GitLabApi) GetFile(owner, repoName, filePid, ref string) (*File, error) {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	fileContent, _, err := gl.Client.RepositoryFiles.GetFile(pid, filePid, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}

	name := filePid
	index := strings.LastIndex(filePid, "/")
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

func (gl *GitLabApi) PushFilesToRepository(files map[string][]byte, owner, repoName, branch, commitMessage string) error {
	// Implement PushFilesToRepository for GitLab
	for filePath, content := range files {
		contentStr := string(content)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String(commitMessage),
			Content:       &contentStr,
		}

		_, _, err := gl.Client.RepositoryFiles.CreateFile(fmt.Sprintf("%s/%s", owner, repoName), filePath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gl *GitLabApi) GetRepositoryArchive(owner, repoName, format, ref string) ([]byte, error) {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	archiveOptions := &gitlab.ArchiveOptions{
		Format: &format, // Choose the archive format
	}

	fileData, _, err := gl.Client.Repositories.Archive(pid, archiveOptions)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (gl *GitLabApi) GetLatestCommitHash(owner, repoName, filePid, ref string) (string, error) {
	pid := fmt.Sprintf("%s/%s", owner, repoName)
	fileContent, _, err := gl.Client.RepositoryFiles.GetFile(pid, filePid, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return "", err
	}
	return fileContent.LastCommitID, nil
}
