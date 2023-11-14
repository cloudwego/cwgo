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

package meta

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudwego/cwgo/pkg/consts"
	gv "github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Version           string `yaml:"cwgo version"`      // cwgo version
	CommandType       string `yaml:"commandType"`       // client or server
	CommunicationType string `yaml:"communicationType"` // http or rpc
	Registry          string `yaml:"registry"`
	Resolver          string `yaml:"resolver"`
}

var GoVersion *gv.Version

func init() {
	// valid by unit test already, so no need to check error
	GoVersion, _ = gv.NewVersion(Version)
}

func (manifest *Manifest) InitAndValidate(dir string) error {
	m, err := loadConfigFile(filepath.Join(dir, consts.ManifestFile))
	if err != nil {
		return fmt.Errorf("can not load \".cwgo\", err: %v", err)
	}

	if len(m.Version) == 0 {
		return fmt.Errorf("can not get cwgo version form \".cwgo\", current project doesn't belong to cwgo")
	}

	*manifest = *m
	_, err = gv.NewVersion(manifest.Version)
	if err != nil {
		return fmt.Errorf("invalid cwgo version in \".cwgo\", err: %v", err)
	}

	return nil
}

func (manifest *Manifest) string() string {
	conf, _ := yaml.Marshal(*manifest)

	return consts.CwgoTitle + "\n\n" +
		string(conf)
}

func (manifest *Manifest) Persist(dir string) error {
	file := filepath.Join(dir, consts.ManifestFile)
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0o644))
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.WriteString(manifest.string())
	return err
}

// loadConfigFile load config file from path
func loadConfigFile(path string) (*Manifest, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	file = bytes.TrimPrefix(file, []byte(consts.CwgoTitle))
	if err = yaml.Unmarshal(file, &manifest); err != nil {
		return nil, fmt.Errorf("decode \".cwgo\" failed, err: %v", err)
	}
	return &manifest, nil
}
