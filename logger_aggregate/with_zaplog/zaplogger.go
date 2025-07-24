package with_zaplog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapLogger struct {
	logger *zap.Logger
	atom   *zap.AtomicLevel
}

type LogOutputType int

const (
	LogOutputType_Normal LogOutputType = iota
	LogOutputType_Json
)

// GetLogger
// params
// 0: serviceName
// 1: logLevel "debug" or "info" or "warn" or "error" or "panic" or "fatal"
// 2: outputMode "json" or "normal"
// 3: devMode "devMode" - panic log then exit.
func GetLogger(params ...interface{}) *ZapLogger {
	var (
		serviceName string
		loggerLevel zapcore.Level
		outputMode  LogOutputType
		devMode     bool
	)

	if len(params) > 0 {
		serviceName = params[0].(string)
	}

	if len(params) > 1 {
		loggerLevel = params[1].(zapcore.Level)
	}

	if len(params) > 2 {
		outputMode = params[2].(LogOutputType)
	}

	if len(params) > 3 {
		devMode = params[3].(bool)
	}

	atom := zap.NewAtomicLevel()
	atom.SetLevel(loggerLevel)

	cores := make([]zapcore.Core, 0)
	if outputMode == LogOutputType_Json {
		// third core
		jsonCfg := zapcore.EncoderConfig{
			TimeKey:        "serverTime",
			LevelKey:       "level",
			NameKey:        "serviceName",
			CallerKey:      "file",
			MessageKey:     "message",
			StacktraceKey:  "stackTrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		cores = append(cores,
			zapcore.NewCore(zapcore.NewJSONEncoder(jsonCfg), zapcore.AddSync(os.Stdout), atom))

	} else {
		// console core
		consoleCfg := zap.NewDevelopmentEncoderConfig()
		consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleCfg.EncodeCaller = zapcore.ShortCallerEncoder

		cores = append(cores,
			zapcore.NewCore(zapcore.NewConsoleEncoder(consoleCfg), zapcore.AddSync(os.Stdout), atom))
	}

	multiCore := zapcore.NewTee(cores...)

	logger := zap.New(multiCore, zap.AddCaller(), zap.AddCallerSkip(1))
	if devMode {
		logger = logger.WithOptions(zap.Development())
	}

	return &ZapLogger{logger.Named(serviceName), &atom}
}

// SetServiceName support dynamic setting
func (zl *ZapLogger) SetServiceName(serviceName string) {
	zl.logger.Named(serviceName)
}

// SetLoggerLevel support dynamic setting
func (zl *ZapLogger) SetLoggerLevel(l zapcore.Level) {
	zl.atom.SetLevel(l)
}

func (zl *ZapLogger) Fatal(msg string, fields ...interface{}) {
	zl.logger.Fatal(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Panic(msg string, fields ...interface{}) {
	zl.logger.DPanic(msg, transToZapFields(fields...)...)
}

// Critical fixme: don't support Critical
func (zl *ZapLogger) Critical(msg string, fields ...interface{}) {
	zl.logger.Error(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Error(msg string, fields ...interface{}) {
	zl.logger.Error(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Errorf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Error(msg)
}

func (zl *ZapLogger) Warn(msg string, fields ...interface{}) {
	zl.logger.Warn(msg, transToZapFields(fields...)...)
}

// Notice fixme: don't support Notice
func (zl *ZapLogger) Notice(msg string, fields ...interface{}) {
	zl.logger.Warn(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Info(msg string, fields ...interface{}) {
	zl.logger.Info(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Infof(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Info(msg)
}

func (zl *ZapLogger) Debug(msg string, fields ...interface{}) {
	zl.logger.Debug(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Debugf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Debug(msg)
}

func (zl *ZapLogger) Println(msg string, fields ...interface{}) {
	zl.logger.Debug(msg, transToZapFields(fields...)...)
}

func (zl *ZapLogger) Printf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Info(msg)
}

func transToZapFields(param ...interface{}) []zap.Field {
	fields := make([]zap.Field, len(param))
	for idx, p := range param {
		if p, ok := p.(zap.Field); ok {
			fields[idx] = p
			continue
		}

		fields[idx] = zap.Any(fmt.Sprintf("param_%d", idx), p)
	}
	return fields
}
