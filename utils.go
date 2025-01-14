package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/qingni918/utils/zaplogger"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func EncodeGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	// 确保所有数据都被写入
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CalcFuncCostTime(logKey string, f func()) {
	timeBegin := time.Now()
	defer func() {
		log.Printf("calcFuncCostTime processed, logKey: %s, cost time: %s", logKey, time.Since(timeBegin).String())
	}()
	f()
}

type Logger struct {
	*log.Logger
}

func (l *Logger) PrintErr(err error, exit ...bool) {
	if err != nil {
		l.Println(err.Error())
	}

	if len(exit) > 0 && exit[0] == true {
		os.Exit(1)
	}
}

func GetLogger() *Logger {
	lg := log.Default()
	lg.SetFlags(log.Ldate | log.Ltime | log.Lshortfile /* | log.LUTC*/)
	return &Logger{lg}
}

func GetZapLoggerWith(serviceID string, level zapcore.Level, outputMode int, bdInfo string, devMode bool) *zaplogger.ZapLogger {
	logger := zaplogger.NewLogger(serviceID, level, outputMode, devMode)
	return logger
}

func GetCaller() string {
	_, file, line, ok := runtime.Caller(2)
	callerStr := "GetCaller - failed"
	if ok {
		callerStr = fmt.Sprintf("%s:%d", file, line)
	}
	return callerStr
}

func GetCallerWithSkip(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	callerStr := "GetCaller - failed"
	if ok {
		callerStr = fmt.Sprintf("%s:%d", file, line)
	}
	return callerStr
}
