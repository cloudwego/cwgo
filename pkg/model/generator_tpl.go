package model

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"path/filepath"
	"text/template"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"

	"gorm.io/gen"
)

//go:embed model_tpl.txt
var mergedTemplate string

type GenMethodTmpl struct {
	GormOpen string
	gen.Config
	OnlyModel bool
	// UseRawSQL indicates whether to use raw SQL for the database connection
	UseRawSQL      bool
	Tables         []string
	StrategyParams struct {
		ExcludeTables []string
		Type          string
	}
}

func execTmpl(c *config.ModelArgument) ([]byte, error) {
	tpl, err := template.New("merged").Parse(mergedTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse template fail: %w", err)
	}

	absOutPath, _ := filepath.Abs(c.OutPath)
	data := GenMethodTmpl{
		GormOpen:  buildGormOpen(c),
		UseRawSQL: c.DSN == "" && c.SQLDir != "",
		Tables:    c.Tables,
		OnlyModel: c.OnlyModel,
		Config: gen.Config{
			OutPath:           absOutPath,
			OutFile:           c.OutFile,
			ModelPkgPath:      c.ModelPkgName,
			WithUnitTest:      c.WithUnitTest,
			FieldNullable:     c.FieldNullable,
			FieldSignable:     c.FieldSignable,
			FieldWithIndexTag: c.FieldWithIndexTag,
			FieldWithTypeTag:  c.FieldWithTypeTag,
			Mode:              buildGenMode(c.Mode),
		},
		StrategyParams: struct {
			ExcludeTables []string
			Type          string
		}{
			ExcludeTables: c.ExcludeTables,
			Type:          c.Type,
		},
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("execute template fail: %w", err)
	}
	fmtCode, err := format.Source(buf.Bytes())
	return fmtCode, err
}

func buildGormOpen(c *config.ModelArgument) string {
	abs, _ := filepath.Abs(c.SQLDir)
	switch {
	case c.SQLDir != "":
		return fmt.Sprintf(
			"db, err := gorm.Open(rawsql.New(rawsql.Config{FilePath: []string{%q}}))",
			abs,
		)
	case c.DSN != "" && c.Type != "":
		return fmt.Sprintf(
			"db, err := gorm.Open(%s.Open(%q))",
			consts.DataBaseType(c.Type),
			c.DSN,
		)
	default:
		return ""
	}
}
