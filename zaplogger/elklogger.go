package zaplogger

import (
	"encoding/json"
	"fmt"
	circuit "github.com/rubyist/circuitbreaker"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

// fixme: unmarshal as map, to keep fields alive
type ELKLogger struct {
	Host         string `json:"host"`
	ShortMessage string `json:"msg"`
	Level        int    `json:"level"`
	File         string `json:"file"`
	LoggerLevel  string `json:"level"`
	LoggerTime   string `json:"time"`
}

type ELKLoggerWriter struct {
	url       string
	serviceID string

	// post buf
	postBuff string
	// post notice
	bulkTicker *time.Ticker
	bulkSwitch chan byte
	bulkSize   int

	cb *circuit.Breaker
	sync.Mutex
}

func (ew *ELKLoggerWriter) Init() {

	ew.bulkTicker = time.NewTicker(time.Second)
	ew.bulkSwitch = make(chan byte, 5)
	ew.cb = circuit.NewThresholdBreaker(3)
	go func() {
		for {
			select {
			case <-ew.bulkTicker.C:
				ew.bulkSwitch <- 1

			case <-ew.bulkSwitch:
				if err := ew.post(); err != nil {
					if !ew.cb.Tripped() {
						ew.cb.Fail()
						_, fn, ln, _ := runtime.Caller(0)
						log.Printf("%s:%d %s\n", fn, ln+1, err.Error())
					}

				} else if ew.cb.Tripped() {
					ew.cb.Reset()
				}
			}
		}
	}()

	events := ew.cb.Subscribe()
	go func() {
		for {
			select {
			case e := <-events:
				if e == circuit.BreakerTripped {
					log.Printf("ELKLoggerWriter circuit-breaker event: BreakerTripped")
				} else if e == circuit.BreakerReset {
					log.Printf("ELKLoggerWriter circuit-breaker event: BreakerReset")
				} else if e == circuit.BreakerFail {
					log.Printf("ELKLoggerWriter circuit-breaker event: BreakerFail")
				} else if e == circuit.BreakerReady {
					log.Printf("ELKLoggerWriter circuit-breaker event: BreakerReady")
				}
			}
		}
	}()
}

func (ew *ELKLoggerWriter) Write(logDetail []byte) (int, error) {

	if ew.cb.Tripped() {
		//return 0, errors.New("ELK tripped")
		return 0, nil
	}

	el := make(map[string]interface{}, 0)
	err := json.Unmarshal(logDetail, &el)
	if err != nil {
		return 0, err
	}
	logDetail, _ = json.Marshal(&el)
	strLogDetail := string(logDetail)

	// 添加元json，用于bulk接口
	strLogDetail = fmt.Sprintf("{\"index\":{\"_index\":\"%s\",\"_type\":\"log\"}}\n%s\n", ew.serviceID, strLogDetail)
	//fmt.Println(ew.url, string(strLogDetail))

	ew.Lock()
	ew.postBuff += strLogDetail
	ew.bulkSize += 1
	if ew.bulkSize == 200 {
		ew.bulkSwitch <- 1
	}
	ew.Unlock()

	return 0, nil
}

func (*ELKLoggerWriter) getJsonEncoderConfig() zapcore.EncoderConfig {

	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "host",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (ew *ELKLoggerWriter) post() error {

	ew.Lock()
	defer ew.Unlock()
	if ew.postBuff == "" {
		return nil
	}

	body := strings.NewReader(ew.postBuff)

	req, err := http.NewRequest("POST", ew.url, body)
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	req.Header = header
	rep, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if rep.StatusCode != 200 {
		_, fn, ln, _ := runtime.Caller(0)
		log.Printf("%s:%d %s\n", fn, ln, rep.Status)
		log.Println(ew.postBuff, "\n", ew.url)
		return nil
	}
	ew.resetPostBuff()
	rep.Body.Close()
	return nil
}

func (ew *ELKLoggerWriter) resetPostBuff() {
	ew.postBuff = ""
	ew.bulkSize = 0
}
