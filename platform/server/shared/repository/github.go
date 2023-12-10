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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/google/go-github/v56/github"
	"go.uber.org/zap"
)

const (
	githubURLPrefix  = "https://github.com/"
	regGitHubURL     = `https://github\.com/([^\/]+)\/([^\/]+)\/blob\/([^\/]+)\/(.+)`
	regGithubRepoURL = `https://github\.com/([^\/]+)\/([^\/]+)`
)

type GitHubApi struct {
	client              *github.Client
	repoApiDomain       string
	token               string
	tokenOwner          string
	repoOwner           string
	repoName            string
	tokenExpirationTime time.Time
}

func NewGitHubApi(token, repoOwner, repoName string) (*GitHubApi, error) {
	client, err := utils.NewGithubClient(token)
	if err != nil {
		return nil, err
	}

	// get token info
	tokenOwner, tokenExpirationTime, err := utils.GetGitHubTokenInfo(client)
	if err != nil {
		return nil, err
	}

	// check token has certain repo permission
	isValid, err := utils.ValidateTokenForRepoGitHub(client, repoOwner, repoName)
	if err != nil {
		logger.Logger.Error("validate token for repo failed", zap.Error(err))
		return nil, err
	}

	if !isValid {
		return nil, consts.ErrTokenInvalid
	}

	return &GitHubApi{
		client:              client,
		repoApiDomain:       consts.GitHubDomain,
		token:               token,
		tokenOwner:          tokenOwner,
		repoOwner:           repoOwner,
		repoName:            repoName,
		tokenExpirationTime: tokenExpirationTime,
	}, nil
}

func (a *GitHubApi) GetRepoApiDomain() (domain string) {
	return a.repoApiDomain
}

func (a *GitHubApi) GetToken() (token string) {
	return a.token
}

func (a *GitHubApi) GetTokenOwner() (tokenOwner string) {
	return a.tokenOwner
}

func (a *GitHubApi) GetRepoOwner() (repoOwner string) {
	return a.repoOwner
}

func (a *GitHubApi) GetRepoName() (repoName string) {
	return a.repoName
}

func (a *GitHubApi) CheckTokenIfExpired() bool {
	return a.tokenExpirationTime.Before(time.Now())
}

func (a *GitHubApi) GetRepoDefaultBranch() (string, error) {
	repo, _, err := a.client.Repositories.Get(context.TODO(), a.repoOwner, a.repoName)
	if err != nil {
		return "", err
	}

	return *repo.DefaultBranch, nil
}

func (a *GitHubApi) ValidateRepoBranch(branch string) (bool, error) {
	branchesRes, _, err := a.client.Repositories.ListBranches(context.Background(), a.repoOwner, a.repoName, &github.BranchListOptions{})
	if err != nil {
		return false, err
	}

	for _, branchRes := range branchesRes {
		if *branchRes.Name == branch {
			return true, nil
		}
	}

	return false, nil
}

func (a *GitHubApi) GetRepoBranches() ([]string, error) {
	branchesRes, _, err := a.client.Repositories.ListBranches(context.Background(), a.repoOwner, a.repoName, &github.BranchListOptions{})
	if err != nil {
		return nil, err
	}

	branches := make([]string, len(branchesRes))

	for i, branchRes := range branchesRes {
		branches[i] = *branchRes.Name
	}

	return branches, nil
}

func (a *GitHubApi) ParseFileUrl(url string) (filePid, owner, repoName string, err error) {
	return utils.ParseRepoFileUrl(consts.RepositoryTypeNumGithub, url)
}

func (a *GitHubApi) GetFile(owner, repoName, filePath, ref string) (*File, error) {
	// prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// download the file content from the GitHub repository.
	fileContent, _, err := a.client.Repositories.DownloadContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// read the file content into a byte slice.
	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}

	// create a File struct with the file name and content.
	return &File{
		Name:    filePath,
		Content: content,
	}, nil
}

func (a *GitHubApi) PushFilesToRepository(files map[string][]byte, owner, repoName, branch, commitMessage string) error {
	// get a reference to the default branch
	ref, _, err := a.client.Git.GetRef(context.Background(), owner, repoName, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	// obtain the tree for the default branch
	baseTree, _, err := a.client.Git.GetTree(context.Background(), owner, repoName, *ref.Object.SHA, false)
	if err != nil {
		return err
	}

	// create a new Tree object for the file to be pushed
	var treeEntries []*github.TreeEntry
	for filePath, content := range files {
		treeEntries = append(treeEntries, &github.TreeEntry{
			Path:    github.String(filePath),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		})
	}

	// add a new file to the tree of the default branch
	treeEntries = append(treeEntries, baseTree.Entries...)

	newTree, _, err := a.client.Git.CreateTree(context.Background(), owner, repoName, *ref.Object.SHA, treeEntries)
	if err != nil {
		return err
	}

	// create a new commit object, using the new tree as its foundation
	newCommit, _, err := a.client.Git.CreateCommit(
		context.Background(),
		owner,
		repoName,
		&github.Commit{
			Message: github.String(commitMessage),
			Tree:    newTree,
			Parents: []*github.Commit{{SHA: ref.Object.SHA}},
		},
		&github.CreateCommitOptions{},
	)
	if err != nil {
		return err
	}

	// update branch references to point to new submissions
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
	// prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// define the format for the archive (e.g., "tarball").
	format := "tarball"

	// get the archive link from the GitHub repository.
	archiveLink, _, err := a.client.Repositories.GetArchiveLink(
		context.Background(),
		owner,
		repoName,
		github.ArchiveFormat(format),
		opts,
		3,
	)
	if err != nil {
		return nil, err
	}

	// fetch the archive data from the obtained link.
	resp, err := http.Get(archiveLink.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check if the HTTP status indicates a successful fetch.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch archive: %s", resp.Status)
	}

	// read the archive data into a byte slice.
	archiveData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return archiveData, nil
}

func (a *GitHubApi) GetLatestCommitHash(owner, repoName, filePath, ref string) (string, error) {
	// prepare options with the desired Git reference.
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// get the contents of the specified file from the GitHub repository.
	fileContent, _, _, err := a.client.Repositories.GetContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return "", err
	}

	// extract and return the SHA (commit hash) associated with the file.
	return *fileContent.SHA, nil
}

func (a *GitHubApi) DeleteDirs(owner, repoName string, folderPaths ...string) error {
	for _, folderPath := range folderPaths {
		// define the file path for a .gitkeep file within the folder.
		filePath := fmt.Sprintf("%s/%s", folderPath, ".gitkeep")

		// configure options for committing the delete operation.
		commitOpts := &github.RepositoryContentFileOptions{
			Message: github.String(fmt.Sprintf("Delete folder %s", folderPath)),
			Branch:  github.String("main"), // Set the branch where the folder deletion should occur
		}

		// attempt to delete the .gitkeep file, effectively removing the folder.
		_, _, err := a.client.Repositories.DeleteFile(context.Background(), owner, repoName, filePath, commitOpts)

		// check if an error occurred during the delete operation.
		if err != nil && !utils.IsFileNotFoundError(err) {
			return err
		}
	}

	return nil
}

func (a *GitHubApi) AutoCreateRepository(owner, repoName string, isPrivate bool) (string, error) {
	ctx := context.Background()
	// new repository's URL
	newRepoURL := githubURLPrefix + owner + "/" + repoName
	_, _, err := a.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		// if the error is caused by the inability to find a repository with the name, create the repository
		if _, ok := err.(*github.ErrorResponse); ok {
			newRepo := &github.Repository{
				Name:        github.String(repoName),
				Private:     &isPrivate,
				Description: github.String("generate by cwgo"),
				AutoInit:    github.Bool(true),
			}

			_, _, err := a.client.Repositories.Create(ctx, "", newRepo)
			if err != nil {
				return "", err
			}

			return newRepoURL, nil
		}
		return "", err
	}

	return newRepoURL, nil
}

func (a *GitHubApi) GetRepositoryPrivacy(owner, repoName string) (bool, error) {
	ctx := context.Background()
	repo, _, err := a.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return false, err
	}

	return repo.GetPrivate(), err
}
