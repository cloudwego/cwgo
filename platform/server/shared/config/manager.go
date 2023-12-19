/*
 * Copyright 2023 CloudWeGo Authors
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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/cwgo/platform/server/shared/config/app"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/api"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/config/internal/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/config/store"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	ServerType         consts.ServerType
	ServerMode         consts.ServerMode
	ServiceId          string
	Config             Config
	ApiConfigManager   *api.ConfigManager
	AgentConfigManager *agent.ConfigManager
}

type Config struct {
	App      app.Config      `mapstructure:"app"`
	Logger   logger.Config   `mapstructure:"logger"`
	Registry registry.Config `mapstructure:"registry"`
	Store    store.Config    `mapstructure:"store"`
	Api      api.Config      `mapstructure:"api"`
	Agent    agent.Config    `mapstructure:"agent"`
}

func (conf *Config) SetUp() {
	conf.App.SetUp()
	conf.Logger.SetUp()
	conf.Registry.SetUp()
	conf.Store.SetUp()
	conf.Api.SetUp()
	conf.Agent.SetUp()
}

var manager *Manager

type FileConfig struct {
	Path string
}

func InitManager(serverType consts.ServerType, serverMode consts.ServerMode, configType consts.ConfigType, metadata ...interface{}) error {
	var config Config

	switch configType {
	case consts.ConfigTypeNumFile:
		var configPath string

		if metadata != nil {
			if fileConfig, ok := metadata[0].(FileConfig); ok {
				configPath = fileConfig.Path
			}
		}

		configPath = filepath.ToSlash(filepath.Join(configPath, fmt.Sprintf("config-%s.yaml", consts.ServerModeMapToStr[serverMode])))

		fmt.Printf("get config path: %s\n", configPath)

		v := viper.New()
		v.SetConfigType("yaml")
		v.SetConfigFile(configPath)
		err := v.ReadInConfig()
		if err != nil {
			panic(fmt.Sprintf("get config file failed, err: %v", err))
		}

		if err := v.Unmarshal(&config); err != nil {
			return fmt.Errorf("unmarshal Config failed, err: %v", err)
		}

	case consts.ConfigTypeNumApollo:
		// TODO: to be implemented
		panic("to be implemented")
	default:

	}

	config.SetUp()

	var err error

	// init consts in config
	consts.ProxyUrl = config.App.ProxyUrl

	if config.App.Timezone == "" {
		consts.TimeZone = time.Local
	} else {
		consts.TimeZone, err = time.LoadLocation(config.App.Timezone)
		if err != nil {
			return err
		}
	}

	// get service id
	var serviceId string
	_, err = os.Stat(consts.AgentMetadataFile)
	if os.IsNotExist(err) {
		// agent file not exist
		// generate a new service id
		serviceId, err = utils.NewServiceId()
		if err != nil {
			return err
		}
	} else {
		// use exist service id
		yamlFileBytes, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			panic(fmt.Sprintf("read agent metadata file failed, err: %v", err))
		}

		var agentMetadata agent.Metadata
		err = yaml.Unmarshal(yamlFileBytes, &agentMetadata)
		if err != nil {
			panic(fmt.Sprintf("unmarshal agent metadata file failed, err: %v", err))
		}

		serviceId = agentMetadata.ServiceId
	}

	switch serverType {
	case consts.ServerTypeNumApi:
		manager = &Manager{
			ServerType:       serverType,
			ServerMode:       serverMode,
			ServiceId:        serviceId,
			Config:           config,
			ApiConfigManager: api.NewConfigManager(config.Api, config.Registry, config.Store, serviceId),
		}
	case consts.ServerTypeNumAgent:
		manager = &Manager{
			ServerType:         serverType,
			ServerMode:         serverMode,
			ServiceId:          serviceId,
			Config:             config,
			AgentConfigManager: agent.NewConfigManager(config.Agent, config.Registry, config.Store, serviceId),
		}
	}

	return nil
}

func GetManager() *Manager {
	if manager == nil {
		panic("config manager not initialized")
	}

	return manager
}
