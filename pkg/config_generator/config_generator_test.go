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
	"testing"
)

func TestConvertToGoStructName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"config_center.yml", "ConfigCenter"},
		{"config-center_test", "ConfigCenterTest"},
		{"dss/config-center", "ConfigCenter"},
		{"config-center-test", "ConfigCenterTest"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToGoStructName(tt.name); got != tt.want {
				t.Errorf("convertToGoStructName() = %v, want %v", got, tt.want)
			}
		})
	}
}
