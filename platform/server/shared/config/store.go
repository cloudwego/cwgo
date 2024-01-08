/*
 *
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
 *
 */

package config

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"

	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreConfig struct {
	Type  string `mapstructure:"type"`
	Mysql Mysql  `mapstructure:"mysql"`
	Redis Redis  `mapstructure:"redis"`
}

type Mysql struct {
	Addr     string `mapstructure:"addr"`
	Port     string `mapstructure:"port"`
	Db       string `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`
}

func (m Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		m.Username,
		m.Password,
		m.Addr,
		m.Port,
		m.Db,
		m.Charset,
		url.PathEscape(consts.TimeZone.String()),
	)
}

func (conf *StoreConfig) NewMysqlDB() (*gorm.DB, error) {
	log.Info("connecting mysql", zap.Reflect("dsn", conf.Mysql.GetDsn()))

	gormLogger, err := log.GetGormZapWriter(log.GetGormLoggerConfig())
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(conf.Mysql.GetDsn()), &gorm.Config{
		Logger:      gormLogger,
		PrepareStmt: true,
	})
	if err != nil {
		log.Error("connect mysql failed", zap.Error(err))
		return nil, err
	}

	return db, err
}

type Redis struct {
	Type       string          `mapstructure:"type"`
	StandAlone RedisStandAlone `mapstructure:"standalone"`
	Cluster    RedisCluster    `mapstructure:"cluster"`
}

type RedisStandAlone struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type RedisCluster struct {
	MasterNum int `mapstructure:"masterNum"`
	Addrs     []*struct {
		Ip   string `mapstructure:"ip"`
		Port string `mapstructure:"port"`
	} `mapstructure:"addrs"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (conf *StoreConfig) NewRedisClient() (redis.UniversalClient, error) {
	var rdb redis.UniversalClient

	if conf.Redis.Type == "standalone" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     conf.Redis.StandAlone.Addr,
			Username: conf.Redis.StandAlone.Username,
			Password: conf.Redis.StandAlone.Password,
			DB:       conf.Redis.StandAlone.Db,
		})
	} else if conf.Redis.Type == "cluster" || conf.Redis.Type == "" {
		addrs := make([]string, len(conf.Redis.Cluster.Addrs))
		for i, addr := range conf.Redis.Cluster.Addrs {
			addrs[i] = fmt.Sprintf("%s:%s", addr.Ip, addr.Port)
		}

		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Username: conf.Redis.Cluster.Username,
			Password: conf.Redis.Cluster.Password,
		})
	} else {
		return nil, errors.New("invalid redis type")
	}

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func (conf *StoreConfig) Init() {
	conf.Type = consts.StoreTypeMysql
}

func (conf *StoreConfig) GetStoreType() consts.StoreType {
	return consts.StoreTypeMapToNum[conf.Type]
}
