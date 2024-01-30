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

package tpl

import (
	"embed"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
)

//go:embed kitex
var kitexTpl embed.FS

//go:embed hertz
var hertzTpl embed.FS

var (
	KitexDir = path.Join(os.TempDir(), consts.Kitex)
	HertzDir = path.Join(os.TempDir(), consts.Hertz)
)

func Init() {
	os.RemoveAll(KitexDir)
	os.RemoveAll(HertzDir)
	os.Mkdir(KitexDir, 0o755)
	os.Mkdir(HertzDir, 0o755)
	initDir(kitexTpl, consts.Kitex, KitexDir)
	initDir(hertzTpl, consts.Hertz, HertzDir)
}

func initDir(fs embed.FS, srcDir, dstDir string) {
	files, err := fs.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {

		newDstPath := path.Join(dstDir, f.Name())
		newSrcPath := path.Join(srcDir, f.Name())

		if f.IsDir() {
			os.Mkdir(newDstPath, 0o755)
			initDir(fs, newSrcPath, newDstPath)
			continue
		}

		content, err := fs.ReadFile(newSrcPath)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile(newDstPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			panic(err)
		}
		file.Write(content)
		file.Close()
	}
}

func RegisterTemplateFunc() {
	for k, f := range sprig.FuncMap() {
		generator.AddTemplateFunc(k, f)
	}
	generator.AddTemplateFunc("ToCamel", func(name string) string {
		name = strings.Replace(name, "_", " ", -1)
		name = strings.Title(name)
		return strings.Replace(name, " ", "", -1)
	})
}
