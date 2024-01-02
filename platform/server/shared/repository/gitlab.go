/*
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
 */

package repository

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

type GitLabApi struct {
	client              *gitlab.Client
	repoApiDomain       string
	token               string
	tokenType           int32
	tokenOwner          string
	tokenOwnerId        int64
	repoOwner           string
	repoName            string
	branch              string
	tokenExpirationTime time.Time
}

func NewGitLabApi(domain, token, repoOwner, repoName, branch string) (*GitLabApi, error) {
	client, err := utils.NewGitlabClient(token, "https://"+domain)
	if err != nil {
		if utils.IsNetworkError(err) {
			return nil, errors.New("client initialization request timeout")
		} else {
			return nil, errors.New("client initialization request unknown error")
		}
	}

	// get token info
	tokenOwner, tokenOwnerId, tokenType, tokenExpirationTime, err := utils.GetGitLabTokenInfo(client)
	if err != nil {
		return nil, err
	}

	// check token has certain repo permission
	isValid, err := utils.ValidateTokenForRepoGitLab(client, repoOwner, repoName)
	if err != nil {
		log.Error("validate token for repo failed", zap.Error(err))
		return nil, err
	}

	if !isValid {
		return nil, consts.ErrTokenInvalid
	}

	return &GitLabApi{
		client:              client,
		repoApiDomain:       domain,
		token:               token,
		tokenType:           tokenType,
		tokenOwner:          tokenOwner,
		tokenOwnerId:        tokenOwnerId,
		repoOwner:           repoOwner,
		repoName:            repoName,
		branch:              branch,
		tokenExpirationTime: tokenExpirationTime,
	}, nil
}

func (a *GitLabApi) GetProjectPid() string {
	return fmt.Sprintf("%s/%s", a.repoOwner, a.repoName)
}

func (a *GitLabApi) GetRepoApiDomain() (domain string) {
	return a.repoApiDomain
}

func (a *GitLabApi) GetToken() (token string) {
	return a.token
}

func (a *GitLabApi) GetTokenOwner() (tokenOwner string) {
	return a.tokenOwner
}

func (a *GitLabApi) GetRepoOwner() (repoOwner string) {
	return a.repoOwner
}

func (a *GitLabApi) GetRepoName() (repoName string) {
	return a.repoName
}

func (a *GitLabApi) GetBranch() (branch string) {
	return a.branch
}

func (a *GitLabApi) UpdateBranch(branch string) {
	a.branch = branch
}

func (a *GitLabApi) CheckTokenIfExpired() bool {
	return a.tokenExpirationTime.Before(time.Now())
}

func (a *GitLabApi) GetRepoDefaultBranch() (string, error) {
	project, _, err := a.client.Projects.GetProject(a.GetProjectPid(), &gitlab.GetProjectOptions{})
	if err != nil {
		return "", err
	}

	return project.DefaultBranch, nil
}

func (a *GitLabApi) ValidateRepoBranch(branch string) (bool, error) {
	branchesRes, _, err := a.client.Branches.ListBranches(a.GetProjectPid(), &gitlab.ListBranchesOptions{})
	if err != nil {
		return false, err
	}

	for _, branchRes := range branchesRes {
		if branchRes.Name == branch {
			return true, nil
		}
	}

	return false, nil
}

func (a *GitLabApi) GetRepoBranches() ([]string, error) {
	branchesRes, _, err := a.client.Branches.ListBranches(a.GetProjectPid(), &gitlab.ListBranchesOptions{})
	if err != nil {
		return nil, err
	}

	branches := make([]string, len(branchesRes))

	for i, branchRes := range branchesRes {
		branches[i] = branchRes.Name
	}

	return branches, nil
}

func (a *GitLabApi) ParseFileUrl(url string) (filePid, owner, repoName string, err error) {
	return utils.ParseRepoFileUrl(consts.RepositoryTypeNumGitLab, url)
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
	repoPid := fmt.Sprintf("%s/%s", owner, repoName)

	// create temp branch
	tempBranch := fmt.Sprintf("cwgo-temp-%d", time.Now().UnixNano())
	_, _, err := a.client.Branches.CreateBranch(repoPid, &gitlab.CreateBranchOptions{
		Branch: &tempBranch,
		Ref:    &branch,
	})
	if err != nil {
		log.Warn("create repo temp branch failed",
			zap.Error(err),
			zap.String("repo_pid", repoPid),
			zap.String("base_branch", branch),
			zap.String("temp_branch", tempBranch),
		)
		return err
	}

	// delete all file in temp branch
	err = a.DeleteAllFiles(owner, repoName, tempBranch)
	if err != nil {
		if !strings.Contains(err.Error(), "doesn't exist") {
			log.Warn("delete all file in temp branch failed",
				zap.Error(err),
				zap.String("repo_pid", repoPid),
				zap.String("branch", tempBranch),
			)
			return err
		}
		// ignore err
		// because if in this case, there is no generate code in repo
	}

	opts := &gitlab.CreateCommitOptions{
		Branch:        &tempBranch,
		CommitMessage: &commitMessage,
	}

	commitActionOptions := make([]*gitlab.CommitActionOptions, len(files))

	i := 0
	for filePath, content := range files {
		commitActionOptions[i] = &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileCreate),
			FilePath: gitlab.String(filePath),
			Content:  gitlab.String(string(content)),
		}
		i++
	}

	opts.Actions = commitActionOptions

	// create commit(push all file) into temp branch
	_, _, err = a.client.Commits.CreateCommit(repoPid, opts)
	if err != nil {
		log.Warn("create commit into temp branch failed",
			zap.Error(err),
			zap.String("repo_pid", repoPid),
			zap.String("branch", tempBranch),
		)
		return err
	}

	// create merge request
	createMergeRequestRes, _, err := a.client.MergeRequests.CreateMergeRequest(repoPid, &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String("cwgo mr"),
		SourceBranch: &tempBranch,
		TargetBranch: &branch,
	})
	if err != nil {
		log.Warn("create merge request from temp branch to source branch failed",
			zap.Error(err),
			zap.String("repo_pid", repoPid),
			zap.String("temp_branch", tempBranch),
			zap.String("source_branch", branch),
		)
		return err
	}

	// approve merge request
	err = retry.Do(
		// if merge a request immediately after create a mr
		// gitlab will return 405
		// so here set a retry mechanism
		func() error {
			_, _, err = a.client.MergeRequests.AcceptMergeRequest(repoPid, createMergeRequestRes.IID, &gitlab.AcceptMergeRequestOptions{
				MergeCommitMessage:        gitlab.String(commitMessage),
				SquashCommitMessage:       gitlab.String(commitMessage),
				MergeWhenPipelineSucceeds: gitlab.Bool(false),
				Squash:                    gitlab.Bool(true),
				ShouldRemoveSourceBranch:  gitlab.Bool(true),
			})
			if err != nil {
				return err
			}

			return nil
		},
		retry.Attempts(5),
		retry.Delay(3*time.Second),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		log.Warn("approve merge request failed",
			zap.Error(err),
			zap.Int("mr_iid", createMergeRequestRes.IID),
			zap.String("repo_pid", repoPid),
			zap.String("temp_branch", tempBranch),
			zap.String("source_branch", branch),
		)

		// if not able to accept a mr
		// then delete the temp branch
		go func() {
			_, err = a.client.Branches.DeleteBranch(repoPid, tempBranch)
			if err != nil {
				log.Warn("delete temp branch failed",
					zap.Error(err),
					zap.String("repo_id", repoPid),
					zap.String("temp_branch", tempBranch),
				)
			}
		}()
		return err
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
		SHA:    &ref,
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

func (a *GitLabApi) DeleteFiles(owner, repoName, branch string, filePaths ...string) error {
	// generate the project ID by combining owner and repoName
	repoPid := fmt.Sprintf("%s/%s", owner, repoName)

	opts := &gitlab.CreateCommitOptions{
		Branch:        &branch,
		CommitMessage: gitlab.String("delete file"),
	}

	commitActionOptions := make([]*gitlab.CommitActionOptions, len(filePaths))

	for i, path := range filePaths {
		commitActionOptions[i] = &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileDelete),
			FilePath: gitlab.String(path),
		}
	}

	opts.Actions = commitActionOptions

	// create commit(push all file) into temp branch
	_, _, err := a.client.Commits.CreateCommit(repoPid, opts)
	if err != nil {
		if !strings.Contains(err.Error(), "doesn't exist") {
			log.Warn("create commit into branch failed",
				zap.Error(err),
				zap.String("repo_pid", repoPid),
				zap.String("branch", branch),
			)
		}
		return err
	}

	return nil
}

func (a *GitLabApi) DeleteAllFiles(owner, repoName, branch string) error {
	// generate the project ID by combining owner and repoName
	repoPid := fmt.Sprintf("%s/%s", owner, repoName)

	// get repo file tree
	tree, _, err := a.client.Repositories.ListTree(repoPid, &gitlab.ListTreeOptions{
		Ref: gitlab.String(branch),
	})
	if err != nil {
		return err
	}

	opts := &gitlab.CreateCommitOptions{
		Branch:        &branch,
		CommitMessage: gitlab.String("delete file"),
	}

	commitActionOptions := make([]*gitlab.CommitActionOptions, len(tree))

	for i, node := range tree {
		commitActionOptions[i] = &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileDelete),
			FilePath: gitlab.String(node.Path),
		}
	}

	opts.Actions = commitActionOptions

	// create commit(push all file) into temp branch
	_, _, err = a.client.Commits.CreateCommit(repoPid, opts)
	if err != nil {
		if !strings.Contains(err.Error(), "doesn't exist") {
			log.Warn("create commit into branch failed",
				zap.Error(err),
				zap.String("repo_pid", repoPid),
				zap.String("branch", branch),
			)
		}
		return err
	}

	return nil
}

func (a *GitLabApi) AutoCreateRepository(owner, repoName string, isPrivate bool) (string, error) {
	// new repository's path in gitlab
	repoPid := owner + "/" + repoName
	repo, _, err := a.client.Projects.GetProject(repoPid, nil)
	if err != nil {
		// if the error is caused by the inability to find a repository with the name, create the repository
		if strings.Contains(err.Error(), "404 Project Not Found") {
			var v gitlab.VisibilityValue
			if isPrivate {
				v = gitlab.PrivateVisibility
			} else {
				v = gitlab.PublicVisibility
			}
			repo, _, err = a.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Name:                         gitlab.String(repoName),
				Visibility:                   &v,
				Description:                  gitlab.String("generate by cwgo"),
				InitializeWithReadme:         gitlab.Bool(true),
				DefaultBranch:                gitlab.String(consts.MainRef),
				AllowMergeOnSkippedPipeline:  gitlab.Bool(true),
				MergePipelinesEnabled:        gitlab.Bool(false),
				RemoveSourceBranchAfterMerge: gitlab.Bool(true),
				// TODO: if repo is org's repo and token is personal token, it will create personal repo
				NamespaceID: gitlab.Int(int(a.tokenOwnerId)),
			})
			if err != nil {
				return "", err
			}

			return repo.WebURL, nil
		}
		return "", err
	}
	return repo.WebURL, nil
}

func (a *GitLabApi) GetRepositoryPrivacy(owner, repoName string) (bool, error) {
	repoPid := owner + "/" + repoName
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
