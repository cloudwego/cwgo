{{ define "GORM_OPEN" }}
    {{ .GormOpen }}
    if err != nil {
        log.Fatal("open db fail: ", err)
    }
{{ end }}

{{ define "GEN_CONFIG" }}
    genConfig := gen.Config{
        OutPath:           {{ printf "%q" .OutPath }},
        OutFile:           {{ printf "%q" .OutFile }},
        ModelPkgPath:      {{ printf "%q" .ModelPkgPath }},
        WithUnitTest:      {{ .WithUnitTest }},
        FieldNullable:     {{ .FieldNullable }},
        FieldSignable:     {{ .FieldSignable }},
        FieldWithIndexTag: {{ .FieldWithIndexTag }},
        FieldWithTypeTag:  {{ .FieldWithTypeTag }},
        Mode:              {{ .Mode }},
    }
{{ end }}

{{ define "TABLE_STRATEGY" }}
    {{if or (gt (len .ExcludeTables) 0) (eq .Type "sqlite") }}
    genConfig.WithTableNameStrategy(func(tableName string) string {
        {{if eq .Type "sqlite" -}}
        if strings.HasPrefix(tableName, "sqlite") {
            return ""
        }
        {{end -}}
        {{if gt (len .ExcludeTables) 0 -}}
        switch tableName {
            {{range $table := .ExcludeTables -}}
            case "{{ $table }}":
                return ""
            {{end -}}
        }
        {{end -}}
        return tableName
    })
    {{end}}
{{ end }}

package main

import (
    "fmt"
    "log"

    "gorm.io/gorm"
    "gorm.io/gen"
    {{if .UseRawSQL }} "gorm.io/rawsql" {{end}}
)

func main() {
    {{ template "GORM_OPEN" . }}
    {{ template "GEN_CONFIG" .Config }}
    {{ template "TABLE_STRATEGY" .StrategyParams }}

    g := gen.NewGenerator(genConfig)
    g.UseDB(db)
	models, err := genModels(g, db, []string{ {{- range $index, $element := .Tables }}
	{{- if $index }}, {{ end }}"{{$element}}"{{- end }} })
	if err != nil {
		log.Fatal("gen models fail: ", err)
	}
	{{if not .OnlyModel}}
		g.ApplyBasic(models...)
	{{end}}
	g.Execute()
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
