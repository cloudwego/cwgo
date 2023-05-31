// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamic

import (
	"flag"
	"strings"

	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
)

var dbConfig = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Select db type",
			Options: []string{"MySQL", "SQLServer", "SQLite", "PostgreSQL"},
		},
		Validate: survey.Required,
		Transform: func(ans interface{}) (newAns interface{}) {
			t := ans.(core.OptionAnswer)
			t.Value = strings.ToLower(t.Value)
			return t
		},
	},
	{
		Name: "DSN",
		Prompt: &survey.Input{
			Message: "Please input db DSN",
		},
		Validate: survey.Required,
	},
}

type ps struct {
	Pass string
}

var pass = []*survey.Question{
	{
		Name: "pass",
		Prompt: &survey.Input{
			Message: "Please input custom param",
		},
	},
}

func parsePass(da *config.ModelArgument, pass string) error {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.StringVar(&da.OutPath, config.OutDir, "biz/dal/query", "")
	f.StringVar(&da.OutFile, config.OutFile, "gen.go", "")
	f.BoolVar(&da.WithUnitTest, config.UnitTest, false, "")
	f.BoolVar(&da.OnlyModel, config.OnlyModel, false, "")
	f.StringVar(&da.ModelPkgName, config.ModelPkgName, "", "")
	f.BoolVar(&da.FieldNullable, config.Nullable, false, "")
	f.BoolVar(&da.FieldSignable, config.Signable, false, "")
	f.BoolVar(&da.FieldWithTypeTag, config.TypeTag, false, "")
	f.BoolVar(&da.FieldWithIndexTag, config.IndexTag, false, "")
	var tables utils.FlagStringSlice
	f.Var(&tables, config.Tables, "")
	if err := f.Parse(utils.StringSliceSpilt([]string{pass})); err != nil {
		return err
	}
	da.Tables = tables
	return nil
}
