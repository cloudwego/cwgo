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

import "fmt"

type Result struct {
	ServiceName           string              `json:"service_name,omitempty"`
	Addr                  string              `json:"addr,omitempty"`
	SubConfigMetadataList []SubConfigMetadata `json:"sub_config_metadata_list,omitempty"`
}

type SubConfigMetadata struct {
	Namespace      string               `json:"namespace,omitempty"`
	ConfigMetadata []ConfigGenerateMeta `json:"config_metadata,omitempty"`
}

// HandleRequest processes the configuration request and returns the result.
func HandleRequest(req *Config) (*Result, error) {
	var addr string
	if req.Addr == nil {
		addr = ""
	} else {
		addr = *req.Addr
	}

	// Create a new result object
	result := &Result{
		ServiceName: req.ServiceName,
		Addr:        addr,
	}

	// Iterate over each sub-config and process it
	for _, subConfig := range req.SubConfigList {
		subConfigMetadata, err := processSubConfig(subConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to process sub-config '%s': %w", subConfig.NameSpace, err)
		}

		// Append processed sub-config metadata to the result
		result.SubConfigMetadataList = append(result.SubConfigMetadataList, subConfigMetadata)
	}

	return result, nil
}

// processSubConfig processes a single sub-config and returns the metadata.
func processSubConfig(subConfig *SubConfig) (SubConfigMetadata, error) {
	metadata := SubConfigMetadata{
		Namespace: subConfig.NameSpace,
	}

	// Iterate over each key-value pair and process it
	for _, configKvPair := range subConfig.ConfigKvPairList {
		fileMetas, err := processConfigKvPair(configKvPair)
		if err != nil {
			return metadata, err
		}

		// Add the processed file metadata to the sub-config metadata
		metadata.ConfigMetadata = append(metadata.ConfigMetadata, fileMetas...)
	}

	return metadata, nil
}

// processConfigKvPair processes a single key-value pair and returns the file metadata.
func processConfigKvPair(configKvPair *ConfigKvPair) ([]ConfigGenerateMeta, error) {
	var result []ConfigGenerateMeta

	key := configKvPair.Key
	content := configKvPair.Value
	desc := configKvPair.Desc
	group := configKvPair.Kind
	valueType := configKvPair.ValueType

	// Check if the value type is either JsonType or YamlType
	switch valueType {
	case ConfigValueType_YamlType, ConfigValueType_JsonType:
		// Convert configuration content into Go structs
		yaml2Go := New(key, desc, group, valueType)
		if _, err := yaml2Go.Convert(convertToGoStructName(key), []byte(content)); err != nil {
			return nil, fmt.Errorf("failed to convert content for key '%s': %w", key, err)
		}

		// Organize the resulting structs
		organizeStructs(yaml2Go.StructsMeta)
		result = append(result, *yaml2Go.StructsMeta)
	default:
		// Store the raw text content for unsupported types
		result = append(result, ConfigGenerateMeta{
			Desc:            desc,
			Kind:            group,
			ConfigValueType: valueType,
			Key:             key,
		})
	}

	return result, nil
}
