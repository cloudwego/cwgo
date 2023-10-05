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

type IRepository interface {
	ParseUrl(url string) (filePid, owner, repoName string, err error)
	GetFile(owner, repoName, filePid, ref string) (*File, error)
	PushFilesToRepository(files map[string][]byte, owner, repoName, branch, commitMessage string) error
	GetRepositoryArchive(owner, repoName, format, ref string) ([]byte, error) //Specify the compressed file format and obtain the byte of the compressed package, gitlab could not specify ref
	GetLatestCommitHash(owner, repoName, filePid, ref string) (string, error) //Get the latest commit hash for the specified file
}

type File struct {
	Name    string
	Content []byte
}
