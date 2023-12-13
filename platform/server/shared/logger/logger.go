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

package logger

import (
	"fmt"
	"os"

	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type Config struct {
	SavePath     string `mapstructure:"savePath"`
	EncoderType  string `mapstructure:"encoderType"`
	EncodeLevel  string `mapstructure:"encodeLevel"`
	EncodeCaller string `mapstructure:"encodeCaller"`
}

func InitLogger(config Config, serverType consts.ServerType, serviceId string, serverMode consts.ServerMode) {
	// create log save dir
	savePath := config.SavePath
	if savePath == "" {
		savePath = fmt.Sprintf("%s-%s/log",
			consts.ServerTypeMapToStr[serverType],
			serviceId,
		)
	} else {
	}

	err := utils.IsNotExistMkDir(savePath)
	if err != nil {
		panic(err)
	}

	switch config.EncoderType {
	case JsonEncoder, ConsoleEncoder:
	default:
		config.EncoderType = ConsoleEncoder
	}

	switch config.EncodeLevel {
	case LowercaseLevelEncoder, LowercaseColorLevelEncoder, CapitalLevelEncoder, CapitalColorLevelEncoder:
	default:
		config.EncodeLevel = CapitalLevelEncoder

	}

	switch config.EncodeCaller {
	case ShortCallerEncoder, FullCallerEncoder:
	default:
		config.EncodeCaller = FullCallerEncoder

	}

	encoder := getEncoder(EncoderOptions{
		EncoderType:  config.EncoderType,
		EncodeLevel:  config.EncodeLevel,
		EncodeCaller: config.EncodeCaller,
	})

	dynamicLevel := zap.NewAtomicLevel()

	switch serverMode {
	case consts.ServerModeNumDev:
		// set current log level to Debug
		dynamicLevel.SetLevel(zap.DebugLevel)

		// debug level
		debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.DebugLevel
		})
		// info level
		infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.InfoLevel
		})
		// warn level
		warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.WarnLevel
		})
		// error level (include error,panic,fatal)
		errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		})

		cores := [...]zapcore.Core{
			zapcore.NewCore(encoder, os.Stdout, dynamicLevel), // console output
			// archived by log level
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/all/server_all.log", savePath)), zapcore.DebugLevel),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/debug/server_debug.log", savePath)), debugPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/info/server_info.log", savePath)), infoPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/warn/server_warn.log", savePath)), warnPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/error/server_error.log", savePath)), errorPriority),
		}
		Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(Logger)

	case consts.ServerModeNumPro:
		dynamicLevel.SetLevel(zap.InfoLevel)

		cores := [...]zapcore.Core{
			zapcore.NewCore(encoder, os.Stdout, dynamicLevel), // console output
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/server.log", savePath)), dynamicLevel),
		}
		Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(Logger)
	default:
		panic("invalid run mode")
	}
}
