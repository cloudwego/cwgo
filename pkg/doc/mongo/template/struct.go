/*
 * Copyright 2024 CloudWeGo Authors
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

package template

import (
	"bytes"

	"github.com/cloudwego/cwgo/pkg/doc/mongo/code"
)

var structTemplate = `{{.Comment}}
type {{.Name}} struct {
{{.StructFields.GetCode}}
}` + "\n"

type StructRender struct {
	Name         string
	Comment      string
	StructFields code.StructFields
}

func (sr *StructRender) RenderObj(buffer *bytes.Buffer) error {
	if err := templateRender(buffer, "structTemplate", structTemplate, sr); err != nil {
		return err
	}
	return nil
}
