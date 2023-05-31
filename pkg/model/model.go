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

	"github.com/cloudwego/cwgo/config"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func Model(c *config.ModelArgument) error {
	dialector := config.OpenTypeFuncMap[config.DataBaseType(c.Type)]
	db, err := gorm.Open(dialector(c.DSN))
	if err != nil {
		return err
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           c.OutPath,
		OutFile:           c.OutFile,
		ModelPkgPath:      c.ModelPkgName,
		WithUnitTest:      c.WithUnitTest,
		FieldNullable:     c.FieldNullable,
		FieldSignable:     c.FieldSignable,
		FieldWithIndexTag: c.FieldWithIndexTag,
	})

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
