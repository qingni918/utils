package zaplogger

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

// devMode - panic log then exit.

// outputMode
// 0: normal log
// 1: json-string log

func NewLogger(serviceID string, loggerLevel zapcore.Level, outputMode int, devMode bool) *ZapLogger {

	atom := zap.NewAtomicLevel()
	atom.SetLevel(loggerLevel)

	cores := make([]zapcore.Core, 0)
	// third core
	if outputMode == 1 {
		jsonCfg := zapcore.EncoderConfig{
			TimeKey:        "serverTime",
			LevelKey:       "level",
			NameKey:        "serviceID",
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
	}

	// console core
	if outputMode == 0 {
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

	return &ZapLogger{logger.Named(serviceID), &atom}
}

// support dynamic setting
func (zl *ZapLogger) SetLoggerLevel(l zapcore.Level) {
	zl.atom.SetLevel(l)
}

// fatal
func (zl *ZapLogger) Fatal(msg string, fields ...zapcore.Field) {
	zl.logger.Fatal(msg, fields...)
}

// panic
func (zl *ZapLogger) Panic(msg string, fields ...zapcore.Field) {
	zl.logger.Panic(msg, fields...)
}

// critical fixme: don't support Critical
func (zl *ZapLogger) Critical(msg string, fields ...zapcore.Field) {
	zl.logger.Error(msg, fields...)
}

// error
func (zl *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	zl.logger.Error(msg, fields...)
}

// errorf
func (zl *ZapLogger) Errorf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Error(msg)
}

// warning
func (zl *ZapLogger) Warning(msg string, fields ...zapcore.Field) {
	zl.logger.Warn(msg, fields...)
}

// notice fixme: don't support Notice
func (zl *ZapLogger) Notice(msg string, fields ...zapcore.Field) {
	zl.logger.Warn(msg, fields...)
}

// info
func (zl *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	zl.logger.Info(msg, fields...)
}

// zk interface
func (zl *ZapLogger) Infof(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Info(msg)
}

// debug
func (zl *ZapLogger) Debug(msg string, fields ...zapcore.Field) {
	zl.logger.Debug(msg, fields...)
}

// zk interface
func (zl *ZapLogger) Printf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	zl.logger.Debug(msg)
}

// zk interface
func (zl *ZapLogger) Println(param ...interface{}) {
	msg := fmt.Sprint(param...)
	zl.logger.Debug(msg)
}
