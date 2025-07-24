package logger_aggregate

import (
	"github.com/qingni918/utils/logger_aggregate/with_log"
	"github.com/qingni918/utils/logger_aggregate/with_zaplog"
)

type Logger interface {
	Println(msg string, fields ...any)
	Printf(format string, fields ...any)
	Debug(msg string, fields ...any)
	Debugf(format string, fields ...any)
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Errorf(format string, fields ...any)
	Warn(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Panic(msg string, fields ...any)
	SetServiceName(serviceName string)
}

type LogType int

const (
	LogType_Log LogType = iota
	LogType_Zaplog
)

var (
	logger  Logger
	logType LogType
)

func initialize(serviceName ...interface{}) {
	if logType == LogType_Log {
		logger = with_log.GetLogger(serviceName...)
	} else if logType == LogType_Zaplog {
		logger = with_zaplog.GetLogger(serviceName...)
	} else {
		logger = with_log.GetLogger(serviceName...)
	}
}

func SetLogType(lt LogType) {
	logType = lt
}

func GetLogger(params ...interface{}) Logger {
	if logger == nil {
		initialize(params...)
	}
	return logger
}
