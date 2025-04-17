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
	"github.com/cloudwego/cwgo/pkg/consts"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var globalArgs = NewArgument()

func GetGlobalArgs() *Argument {
	return globalArgs
}

type Argument struct {
	Verbose bool

	*ServerArgument
	*ClientArgument
	*ModelArgument
	*DocArgument
	*JobArgument
	*ApiArgument
	*DockerArgument
	*KubeArgument
	*FallbackArgument
}

func NewArgument() *Argument {
	return &Argument{
		ServerArgument:   NewServerArgument(),
		ClientArgument:   NewClientArgument(),
		ModelArgument:    NewModelArgument(),
		DocArgument:      NewDocArgument(),
		JobArgument:      NewJobArgument(),
		ApiArgument:      NewApiArgument(),
		FallbackArgument: NewFallbackArgument(),
	}
}

type DialectorFunc func(string) gorm.Dialector

var OpenTypeFuncMap = map[consts.DataBaseType]DialectorFunc{
	consts.MySQL:     mysql.Open,
	consts.SQLServer: sqlserver.Open,
	consts.Sqlite:    sqlite.Open,
	consts.Postgres:  postgres.Open,
}
