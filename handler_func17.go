// +build go1.7

package jsonrpc

import (
	"context"
	"net/http"
)

var (
	// Before runs before invoke a method.
	Before func(context.Context, *Request) *Error
	// After runs after invoke a method.
	After func(context.Context, *Response, *Request)
)

// HandlerFunc provides basic JSON-RPC handling.
func HandlerFunc(w http.ResponseWriter, r *http.Request) {

	rs, batch, err := ParseRequest(r)
	if err != nil {
		SendResponse(w, []Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		return
	}

	c := r.Context()
	resp := make([]Response, len(rs))
	for i := range rs {
		resp[i] = invokeMethod(c, rs[i])
	}

	if err := SendResponse(w, resp, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func invokeMethod(c context.Context, r Request) Response {
	res := NewResponse(r)
	if After != nil {
		defer After(c, &res, &r)
	}
	if Before != nil {
		res.Error = Before(c, &r)
		if res.Error != nil {
			return res
		}
	}
	var h Handler
	h, res.Error = TakeMethod(r)
	if res.Error != nil {
		return res
	}
	res.Result, res.Error = h.ServeJSONRPC(c, r.Params)
	return res
}
