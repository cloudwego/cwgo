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

func Test_HandleRequest(t *testing.T) {
	c := &Config{
		ServiceName: "nacos_config_server",
		SubConfigList: []*SubConfig{
			{
				NameSpace: "public",
				ConfigKvPairList: []*ConfigKvPair{
					{
						Key:       "conf.yaml",
						ValueType: ConfigValueType_YamlType,
						Kind:      "dev",
						Desc:      "dsds",
						Value:     value,
					},
					{
						Key:       "conf.yaml",
						ValueType: ConfigValueType_XmlType,
						Kind:      "dev",
						Value:     value,
					},
				},
			},
		},
	}
	ss, err := HandleRequest(c)
	if err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile("test.json", out, 0o644)
}
