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

package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/urfave/cli/v2"
)

type ModelArgument struct {
	C                 string
	DSN               string   `yaml:"dsn"`
	Type              string   `yaml:"db"`
	Tables            []string `yaml:"tables"`
	OnlyModel         bool     `yaml:"onlyModel"`
	OutPath           string   `yaml:"outPath"`
	OutFile           string   `yaml:"outFile"`
	WithUnitTest      bool     `yaml:"withUnitTest"`
	ModelPkgName      string   `yaml:"modelPkgName"`
	FieldNullable     bool     `yaml:"fieldNullable"`
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"`
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`
}

const (
	DefaultOutPath = "biz/dal/query"
	DefaultOutFile = "gen.go"
)

func NewModelArgument() *ModelArgument {
	return &ModelArgument{
		OutPath: DefaultOutPath,
		OutFile: DefaultOutFile,
	}
}

func (c *ModelArgument) ParseCli(ctx *cli.Context) error {
	configPath := ctx.String(C)
	dsn := ctx.String(DSN)
	tp := strings.ToLower(ctx.String(DBType))
	tables := ctx.StringSlice(Tables)
	onlyModel := ctx.Bool(OnlyModel)
	outPath := ctx.String(OutDir)
	outFile := ctx.String(OutFile)
	withUnitTest := ctx.Bool(UnitTest)
	modelPkgName := ctx.String(ModelPkgName)
	fieldNullable := ctx.Bool(Nullable)
	fieldWithIndexTag := ctx.Bool(IndexTag)
	fieldWithTypeTag := ctx.Bool(TypeTag)
	// priority: command line > config file
	if configPath != "" {
		if configFileParams, err := loadConfigFile(c.C); err == nil && configFileParams != nil {
			c.DSN = configFileParams.DSN
			c.Type = configFileParams.Type
			c.Tables = configFileParams.Tables
			c.OnlyModel = configFileParams.OnlyModel
			if configFileParams.OutPath != "" {
				c.OutPath = configFileParams.OutPath
			}
			if configFileParams.OutFile != "" {
				c.OutFile = configFileParams.OutFile
			}
			c.WithUnitTest = configFileParams.WithUnitTest
			c.ModelPkgName = configFileParams.ModelPkgName
			c.FieldNullable = configFileParams.FieldNullable
			c.FieldWithIndexTag = configFileParams.FieldWithIndexTag
			c.FieldWithTypeTag = configFileParams.FieldWithTypeTag
		} else {
			return err
		}
	}
	if dsn != "" {
		c.DSN = dsn
	}
	if tp != "" {
		c.Type = tp
	}
	if tables != nil {
		c.Tables = tables
	}
	if onlyModel {
		c.OnlyModel = onlyModel
	}
	if outPath != "" {
		c.OutPath = outPath
	}
	if outFile != "" {
		c.OutFile = outFile
	}
	if withUnitTest {
		c.WithUnitTest = withUnitTest
	}
	if modelPkgName != "" {
		c.ModelPkgName = modelPkgName
	}
	if fieldNullable {
		c.FieldNullable = fieldNullable
	}
	if fieldWithIndexTag {
		c.FieldWithIndexTag = fieldWithIndexTag
	}
	if fieldWithTypeTag {
		c.FieldWithTypeTag = fieldWithTypeTag
	}
	return nil
}

// loadConfigFile load config file from path
func loadConfigFile(path string) (*ModelArgument, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint
	var yamlConfig YamlConfig
	if cmdErr := yaml.NewDecoder(file).Decode(&yamlConfig); cmdErr != nil {
		return nil, cmdErr
	}
	return yamlConfig.Database, nil
}

// YamlConfig is yaml config struct
type YamlConfig struct {
	Version  string         `yaml:"version"`  //
	Database *ModelArgument `yaml:"database"` //
}
