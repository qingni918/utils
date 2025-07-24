package httputil

import "net/http"

func ServerAndListen(add string, handler func(*Response, *Request)) error {
	err := http.ListenAndServe(add, Handler(handler))
	if err != nil {
		return err
	}
	return nil
}
