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
	"context"
	"github.com/google/go-github/v53/github"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
)

type IRepository interface {
	GetRepoInfo() (*Repository, error)
	GetFile(filePath, ref string) (*File, error)
	PushFilesToRepository(files map[string][]byte, branch, commitMessage string) error
}

type GitHubApi struct {
	Client   github.Client
	Owner    string
	RepoName string
}

type GitLabApi struct {
	Client    gitlab.Client
	Owner     string
	ProjectID int
}

type File struct {
	Name    string
	Content []byte
}

type Repository struct {
	Name string
	URL  string
}

func (g *GitHubApi) GetRepoInfo() (*Repository, error) {
	repo, _, err := g.Client.Repositories.Get(context.Background(), g.Owner, g.RepoName)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name: repo.GetName(),
		URL:  repo.GetURL(),
	}, nil
}

func (g *GitHubApi) GetFile(filePath, ref string) (*File, error) {
	fileContent, _, err := g.Client.Repositories.DownloadContents(context.Background(), g.Owner, g.RepoName, filePath, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    filePath,
		Content: content,
	}, nil
}

func (g *GitHubApi) PushFilesToRepository(files map[string][]byte, branch, commitMessage string) error {
	// Get a reference to the default branch
	ref, _, err := g.Client.Git.GetRef(context.Background(), g.Owner, g.RepoName, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	// Create a new Tree object for the file to be pushed
	var treeEntries []*github.TreeEntry
	for filePath, content := range files {
		treeEntries = append(treeEntries, &github.TreeEntry{
			Path:    github.String(filePath),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		})
	}
	newTree, _, err := g.Client.Git.CreateTree(context.Background(), g.Owner, g.RepoName, *ref.Object.SHA, treeEntries)
	if err != nil {
		return err
	}

	// Create a new commit object, using the new tree as its foundation
	newCommit, _, err := g.Client.Git.CreateCommit(context.Background(), g.Owner, g.RepoName, &github.Commit{
		Message: github.String(commitMessage),
		Tree:    newTree,
	})
	if err != nil {
		return err
	}

	// Update branch references to point to new submissions
	_, _, err = g.Client.Git.UpdateRef(context.Background(), g.Owner, g.RepoName, &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA:  newCommit.SHA,
			Type: github.String("commit"),
		},
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (gl *GitLabApi) GetRepoInfo() (*Repository, error) {
	proj, _, err := gl.Client.Projects.GetProject(gl.ProjectID, nil)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name: proj.Name,
		URL:  proj.WebURL,
	}, nil
}

func (gl *GitLabApi) GetFile(filePath, ref string) (*File, error) {
	fileContent, _, err := gl.Client.RepositoryFiles.GetFile(gl.ProjectID, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    filePath,
		Content: []byte(fileContent.Content),
	}, nil
}

func (gl *GitLabApi) PushFilesToRepository(files map[string][]byte, branch, commitMessage string) error {
	// Implement PushFilesToRepository for GitLab
	for filePath, content := range files {
		contentStr := string(content)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String(commitMessage),
			Content:       &contentStr,
		}

		_, _, err := gl.Client.RepositoryFiles.CreateFile(gl.ProjectID, filePath, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
