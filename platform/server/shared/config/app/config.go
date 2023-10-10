/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package app

import "time"

type Config struct {
	SyncAgentServiceInterval string `mapstructure:"syncAgentServiceInterval"`
	SyncRepositoryInterval   string `mapstructure:"syncRepositoryInterval"`
	SyncIdlInterval          string `mapstructure:"syncIdlInterval"`
}

const (
	defaultSyncAgentServiceInterval = 3 * time.Second
	defaultSyncRepositoryInterval   = 10 * time.Second
	defaultSyncIdlInterval          = 10 * time.Second
)

func (c Config) GetSyncAgentServiceInterval() time.Duration {
	duration, err := time.ParseDuration(c.SyncAgentServiceInterval)
	if err == nil {
		return duration
	}

	return defaultSyncAgentServiceInterval
}

func (c Config) GetSyncRepositoryInterval() time.Duration {
	duration, err := time.ParseDuration(c.SyncIdlInterval)
	if err == nil {
		return duration
	}

	return defaultSyncRepositoryInterval
}

func (c Config) GetSyncIdlInterval() time.Duration {
	duration, err := time.ParseDuration(c.SyncIdlInterval)
	if err == nil {
		return duration
	}

	return defaultSyncIdlInterval
}
