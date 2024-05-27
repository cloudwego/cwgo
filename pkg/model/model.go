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

package model

import (
	"fmt"
	"strings"

	"gorm.io/rawsql"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"

	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func Model(c *config.ModelArgument) error {
	dialector := config.OpenTypeFuncMap[consts.DataBaseType(c.Type)]

	if c.SQLDir != "" {
		db, err = gorm.Open(rawsql.New(rawsql.Config{
			FilePath: []string{c.SQLDir},
		}))
	} else {
		db, err = gorm.Open(dialector(c.DSN))
	}
	if err != nil {
		return err
	}

	genConfig := gen.Config{
		OutPath:           c.OutPath,
		OutFile:           c.OutFile,
		ModelPkgPath:      c.ModelPkgName,
		WithUnitTest:      c.WithUnitTest,
		FieldNullable:     c.FieldNullable,
		FieldSignable:     c.FieldSignable,
		FieldWithIndexTag: c.FieldWithIndexTag,
	}

	if len(c.ExcludeTables) > 0 || c.Type == string(consts.Sqlite) {
		genConfig.WithTableNameStrategy(func(tableName string) (targetTableName string) {
			if c.Type == string(consts.Sqlite) && strings.HasPrefix(tableName, "sqlite") {
				return ""
			}
			if len(c.ExcludeTables) > 0 {
				for _, table := range c.ExcludeTables {
					if tableName == table {
						return ""
					}
				}
			}
			return tableName
		})
	}

	g := gen.NewGenerator(genConfig)

	g.UseDB(db)

	models, err := genModels(g, db, c.Tables)
	if err != nil {
		return err
	}

	if !c.OnlyModel {
		g.ApplyBasic(models...)
	}

	g.Execute()
	return nil
}

func genModels(g *gen.Generator, db *gorm.DB, tables []string) (models []interface{}, err error) {
	var tablesNameList []string
	if len(tables) == 0 {
		tablesNameList, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("migrator get all tables fail: %w", err)
		}
	} else {
		tablesNameList = tables
	}

	models = make([]interface{}, len(tablesNameList))
	for i, tableName := range tablesNameList {
		models[i] = g.GenerateModel(tableName)
	}
	return models, nil
}
