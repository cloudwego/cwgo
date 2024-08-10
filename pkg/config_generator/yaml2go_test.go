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

package config_generator

import (
	"encoding/json"
	"os"
	"testing"
)

var value = `
kitex:
  service: p.s.m
  version: 1.0.0
  ports:
    - 8888
    - 8889
`

func TestYaml2Go(t *testing.T) {
	yaml2Go := New("key", "desc", "group", ConfigValueType_YamlType)
	s, err := yaml2Go.Convert("config", []byte(value))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	organizeStructs(yaml2Go.StructsMeta)
	marshal, err := json.Marshal(yaml2Go.StructsMeta)
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile("metadata.json", marshal, 0o644)
}
