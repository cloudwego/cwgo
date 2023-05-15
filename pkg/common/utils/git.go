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

package utils

import (
	"net/url"
	"os/exec"
	"strings"
)

func GitClone(gitURL, path string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	c := exec.Command("git", "clone", gitURL)
	c.Dir = path
	return c.Run()
}

func GitPath(gitURL string) (string, error) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}
	p := strings.Split(strings.Trim(u.Path, ""), "/")
	path := p[len(p)-1]
	return path[:len(path)-4], nil
}
