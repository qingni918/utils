package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestPanic(t *testing.T) {
	CalcFuncCostTime("testPanic", func() {
		var err error
		Panic(err)
		err = fmt.Errorf("err")
		Panic(err)
	})
}

func TestEncodeGzip(t *testing.T) {
	str := []byte(`
{"level":"debug","serverTime":"2024-09-05T10:18:48.871+0800","serviceID":"matchshared-10","file":"server/httpserver.go:48","message":"request processed","_pos":"[httpserver.go:1176]","_REQ":"action=share_team_del&caller=match-11&match_phase=rival&match_type=MTYPE_PVE_TRAINING&team_id=118","_REP":"{\"Result\":0,\"Desc\":\"ok\"}"}
`)
	fmt.Println("source len:", len(str))
	encodeStr, err := EncodeGzip(str)
	Panic(err)
	fmt.Println("encode:", len(encodeStr))

	decodeStr, err := DecodeGzip(encodeStr)
	Panic(err)
	fmt.Println("decode:", string(decodeStr))
}

func TestLMHash(t *testing.T) {
	CalcFuncCostTime("testLMHash", func() {
		//log.Println(LMHash("p@ssw0rd"))
		log.Println(LMHash("123456"))
	})
}
