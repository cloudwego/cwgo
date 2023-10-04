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

package logger

import (
	"fmt"
	"github.com/cloudwego/cwgo/platform/server/shared/config"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func InitLogger() {
	// create log save dir
	loggerConfig := config.GetManager().Config.Logger

	savePath := loggerConfig.SavePath
	if savePath == "" {
		savePath = fmt.Sprintf("%s-%s/log",
			consts.ServerTypeMapToStr[config.GetManager().ServerType],
			config.GetManager().ServiceId,
		)
	} else {

	}

	err := utils.IsNotExistMkDir(savePath)
	if err != nil {
		panic(err)
	}

	switch loggerConfig.EncoderType {
	case JsonEncoder, ConsoleEncoder:
	default:
		loggerConfig.EncoderType = ConsoleEncoder
	}

	switch loggerConfig.EncodeLevel {
	case LowercaseLevelEncoder, LowercaseColorLevelEncoder, CapitalLevelEncoder, CapitalColorLevelEncoder:
	default:
		loggerConfig.EncodeLevel = CapitalLevelEncoder

	}

	switch loggerConfig.EncodeCaller {
	case ShortCallerEncoder, FullCallerEncoder:
	default:
		loggerConfig.EncodeCaller = FullCallerEncoder

	}

	encoder := getEncoder(EncoderOptions{
		EncoderType:  loggerConfig.EncoderType,
		EncodeLevel:  loggerConfig.EncodeLevel,
		EncodeCaller: loggerConfig.EncodeCaller,
	})

	dynamicLevel := zap.NewAtomicLevel()

	mode := config.GetManager().ServerMode

	switch mode {
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
