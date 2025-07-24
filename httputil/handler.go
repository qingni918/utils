package httputil

import (
	"net/http"
)

func Handler(f func(*Response, *Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uw := &Response{ResponseWriter: w}
		ur := &Request{Request: r}
		f(uw, ur)
	}
}
