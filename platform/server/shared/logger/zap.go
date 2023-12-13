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
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getEncoder(options EncoderOptions) zapcore.Encoder {
	if options.EncoderType == JsonEncoder {
		return zapcore.NewJSONEncoder(getEncoderConfig(options))
	}

	return zapcore.NewConsoleEncoder(getEncoderConfig(options))
}

func getEncoderConfig(options EncoderOptions) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case options.EncodeLevel == LowercaseLevelEncoder:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case options.EncodeLevel == LowercaseColorLevelEncoder:
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case options.EncodeLevel == CapitalLevelEncoder:
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case options.EncodeLevel == CapitalColorLevelEncoder:
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	if options.EncodeCaller == ShortCallerEncoder {
		config.EncodeCaller = zapcore.ShortCallerEncoder
	}
	return config
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,  // log file save path
		MaxSize:    10,    // log file max size before cut (MB)
		MaxBackups: 50000, // old file max save num
		MaxAge:     1000,  // old file max save days
		Compress:   true,  // is compress old files
		LocalTime:  true,  // is use local time
	}
	return zapcore.AddSync(lumberJackLogger)
}

// CustomTimeEncoder time format
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}
