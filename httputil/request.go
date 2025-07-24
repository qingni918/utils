package httputil

import (
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
	// query 解析
	values       url.Values
	valuesParsed bool

	// body 解析
	body       []byte
	bodyParsed bool
}

func (req *Request) GetRequest() *http.Request {
	return req.Request
}

func (req *Request) ParseQuery() url.Values {

	if req.valuesParsed {
		return req.values
	}

	req.values, _ = url.ParseQuery(req.Request.URL.RawQuery)
	req.valuesParsed = true
	return req.values
}

func (req *Request) ParseBody() []byte {

	if req.bodyParsed {
		return req.body
	}

	defer req.Body.Close()
	req.body, _ = io.ReadAll(req.Body)
	req.bodyParsed = true
	return req.body
}
