package httputil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (rd ResponseData) Bytes() []byte {
	dataBytes, _ := json.Marshal(rd)
	return dataBytes
}

func (resp *Response) Response(code int, message string) {
	resp.Write(ResponseData{Code: code, Msg: message}.Bytes())
}
