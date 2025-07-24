package with_log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// 使用log包装

type Logger struct {
	log          *log.Logger
	serviceName  string
	serviceNamed bool
}

var lg *Logger

func GetLogger(params ...interface{}) *Logger {
	lg = &Logger{log: log.New(os.Stderr, "", 0)}
	if len(params) > 0 {
		lg.SetServiceName(params[0].(string))
	}
	return lg
}

func genPrefix(file string, line int) string {
	if lg.serviceNamed {
		return fmt.Sprintf("%s %s %s:%d: ", time.Now().Format(time.RFC3339), lg.serviceName, file, line)
	}
	return fmt.Sprintf("%s %s:%d: ", time.Now().Format(time.RFC3339), file, line)
}

func (l *Logger) SetServiceName(serviceName string) {
	l.serviceName = serviceName
	// 此处相当于可以取消命名
	l.serviceNamed = serviceName != ""
}

func (l *Logger) Println(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(msg, fmt.Sprint(fields...))
}

func (l *Logger) Printf(format string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(fmt.Sprintf(format, fields...))
}

func (l *Logger) Debug(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(msg, fmt.Sprint(fields...))
}

func (l *Logger) Debugf(format string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(fmt.Sprintf(format, fields...))
}

func (l *Logger) Info(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(msg, fmt.Sprint(fields...))
}

func (l *Logger) Error(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(msg, fmt.Sprint(fields...))
}

func (l *Logger) Errorf(format string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(fmt.Sprintf(format, fields...))
}

func (l *Logger) Warn(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Println(msg, fmt.Sprint(fields...))
}

func (l *Logger) Panic(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Panic(msg, fmt.Sprint(fields...))
}

func (l *Logger) Fatal(msg string, fields ...any) {
	_, file, line, ok := runtime.Caller(1) // 调整为 2 可再多跳一层
	if ok {
		l.log.SetPrefix(genPrefix(file, line))
	}
	l.log.Fatal(msg, fmt.Sprint(fields...))
}
