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
	"errors"
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/google/go-github/v53/github"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type GitHubApi struct {
	client *github.Client
}

func NewGitHubApi(client *github.Client) *GitHubApi {
	return &GitHubApi{
		client: client,
	}
}

const (
	githubURLPrefix = "https://github.com/"
)

func (a *GitHubApi) ParseUrl(url string) (filePid, owner, repoName string, err error) {
	var tempPath string

	// Check if the URL has the GitHub prefix.
	if strings.HasPrefix(url, githubURLPrefix) {
		tempPath = url[len(githubURLPrefix):]

		// Remove any query parameters from the URL.
		lastQuestionMarkIndex := strings.LastIndex(tempPath, "?")
		if lastQuestionMarkIndex != -1 {
			tempPath = tempPath[:lastQuestionMarkIndex]
		}
	} else {
		return "", "", "", errors.New("IDL path format is incorrect; it does not have the expected prefix: " + githubURLPrefix)
	}

	// Define a regular expression to parse the GitHub URL.
	regex := regexp.MustCompile(`([^\/]+)\/([^\/]+)\/blob\/([^\/]+)\/(.+)`)

	// Use the regular expression to extract relevant components from the URL.
	matches := regex.FindStringSubmatch(tempPath)
	if len(matches) != 5 {
		return "", "", "", errors.New("IDL path format is incorrect; unable to parse the GitHub URL")
	}

	// Assign values to the returned variables.
	owner = matches[1]
	repoName = matches[2]
	filePid = matches[4]

	return filePid, owner, repoName, nil
}

func (a *GitHubApi) GetFile(owner, repoName, filePath, ref string) (*File, error) {
	// Prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// Download the file content from the GitHub repository.
	fileContent, _, err := a.client.Repositories.DownloadContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// Read the file content into a byte slice.
	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}

	// Create a File struct with the file name and content.
	return &File{
		Name:    filePath,
		Content: content,
	}, nil
}

func (a *GitHubApi) PushFilesToRepository(files map[string][]byte, owner, repoName, branch, commitMessage string) error {
	// Get a reference to the default branch
	ref, _, err := a.client.Git.GetRef(context.Background(), owner, repoName, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	// Obtain the tree for the default branch
	baseTree, _, err := a.client.Git.GetTree(context.Background(), owner, repoName, *ref.Object.SHA, false)
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

	// Add a new file to the tree of the default branch
	treeEntries = append(treeEntries, baseTree.Entries...)

	newTree, _, err := a.client.Git.CreateTree(context.Background(), owner, repoName, *ref.Object.SHA, treeEntries)
	if err != nil {
		return err
	}

	// Create a new commit object, using the new tree as its foundation
	newCommit, _, err := a.client.Git.CreateCommit(context.Background(), owner, repoName, &github.Commit{
		Message: github.String(commitMessage),
		Tree:    newTree,
		Parents: []*github.Commit{{SHA: ref.Object.SHA}},
	})
	if err != nil {
		return err
	}

	// Update branch references to point to new submissions
	_, _, err = a.client.Git.UpdateRef(context.Background(), owner, repoName, &github.Reference{
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

func (a *GitHubApi) GetRepositoryArchive(owner, repoName, ref string) ([]byte, error) {
	// Prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// Define the format for the archive (e.g., "tarball").
	format := "tarball"

	// Get the archive link from the GitHub repository.
	archiveLink, _, err := a.client.Repositories.GetArchiveLink(context.Background(), owner, repoName, github.ArchiveFormat(format), opts, false)
	if err != nil {
		return nil, err
	}

	// Fetch the archive data from the obtained link.
	resp, err := http.Get(archiveLink.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the HTTP status indicates a successful fetch.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch archive: %s", resp.Status)
	}

	// Read the archive data into a byte slice.
	archiveData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return archiveData, nil
}

func (a *GitHubApi) GetLatestCommitHash(owner, repoName, filePath, ref string) (string, error) {
	// Prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// Get the contents of the specified file from the GitHub repository.
	fileContent, _, _, err := a.client.Repositories.GetContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return "", err
	}

	// Extract and return the SHA (commit hash) associated with the file.
	return *fileContent.SHA, nil
}

func (a *GitHubApi) DeleteDirs(owner, repoName string, folderPaths ...string) error {
	for _, folderPath := range folderPaths {
		// Define the file path for a .gitkeep file within the folder.
		filePath := fmt.Sprintf("%s/%s", folderPath, ".gitkeep")

		// Configure options for committing the delete operation.
		commitOpts := &github.RepositoryContentFileOptions{
			Message: github.String(fmt.Sprintf("Delete folder %s", folderPath)),
			Branch:  github.String("main"), // Set the branch where the folder deletion should occur
		}

		// Attempt to delete the .gitkeep file, effectively removing the folder.
		_, _, err := a.client.Repositories.DeleteFile(context.Background(), owner, repoName, filePath, commitOpts)

		// Check if an error occurred during the delete operation.
		if err != nil && !utils.IsFileNotFoundError(err) {
			return err
		}
	}

	return nil
}
