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
	"github.com/cloudwego/cwgo/pkg/consts"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
)

var dbConfig = []*survey.Question{
	{
		Name: consts.ServiceType,
		Prompt: &survey.Select{
			Message: "Select db type",
			Options: []string{string(consts.MySQL), string(consts.SQLServer), string(consts.Sqlite), string(consts.Postgres)},
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
		Name: consts.Pass,
		Prompt: &survey.Input{
			Message: "Please input custom param",
		},
	},
}

func parsePass(da *config.ModelArgument, pass string) error {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.StringVar(&da.OutPath, consts.OutDir, consts.DefaultDbOutDir, "")
	f.StringVar(&da.OutFile, consts.OutFile, consts.DefaultDbOutFile, "")
	f.BoolVar(&da.WithUnitTest, consts.UnitTest, false, "")
	f.BoolVar(&da.OnlyModel, consts.OnlyModel, false, "")
	f.StringVar(&da.ModelPkgName, consts.ModelPkgName, "", "")
	f.BoolVar(&da.FieldNullable, consts.Nullable, false, "")
	f.BoolVar(&da.FieldSignable, consts.Signable, false, "")
	f.BoolVar(&da.FieldWithTypeTag, consts.TypeTag, false, "")
	f.BoolVar(&da.FieldWithIndexTag, consts.IndexTag, false, "")
	var (
		tables        utils.FlagStringSlice
		excludeTables utils.FlagStringSlice
	)
	f.Var(&tables, consts.Tables, "")
	f.Var(&excludeTables, consts.ExcludeTables, "")
	if err := f.Parse(utils.StringSliceSpilt([]string{pass})); err != nil {
		return err
	}
	da.Tables = tables
	da.ExcludeTables = excludeTables
	return nil
}
