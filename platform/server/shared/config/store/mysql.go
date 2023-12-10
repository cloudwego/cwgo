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

package store

import (
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
)

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

func (c Config) NewMysqlDB() (*gorm.DB, error) {
	logger.Logger.Info("connecting mysql", zap.Reflect("dsn", c.Mysql.GetDsn()))

	gormLogger, err := logger.GetGormZapWriter(logger.GetGormLoggerConfig())
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(c.Mysql.GetDsn()), &gorm.Config{
		Logger:      gormLogger,
		PrepareStmt: true,
	})
	if err != nil {
		logger.Logger.Error("connect mysql failed", zap.Error(err))
		return nil, err
	}

	return db, err
}
