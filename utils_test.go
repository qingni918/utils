package utils

import (
	"fmt"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/qingni918/utils/feishuclient"
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

func TestFeishuClient(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// cli_a7f766e76d60d013
	// I5AyPWY6MAuBx6hhDyAFxgxo7nN8HfYm
	fc := feishuclient.NewClient("cli_a7f766e76d60d013", "I5AyPWY6MAuBx6hhDyAFxgxo7nN8HfYm")
	fc.SetPageSize(50)
	departments := fc.DepartmentChildrenQueryAll("0", "", true)

	type Department struct {
		*larkcontact.Department
		Children []*Department `json:"children,omitempty"`
	}

	depts := make(map[string]*Department, len(departments))

	// 1. 建立索引
	for _, v := range departments {
		depts[*v.OpenDepartmentId] = &Department{
			Department: v,
		}
	}

	// 2. 挂载父子关系
	var rootDepts []*Department
	for _, v := range departments {
		id := *v.OpenDepartmentId
		pid := *v.ParentDepartmentId

		if pid == "0" {
			rootDepts = append(rootDepts, depts[id])
			continue
		}

		if parent, ok := depts[pid]; ok {
			parent.Children = append(parent.Children, depts[id])
		}
	}

	log.Println(len(departments), JsonString(rootDepts))
}
