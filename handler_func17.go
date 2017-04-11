// +build go1.7

package jsonrpc

import (
	"context"
	"net/http"
)

// HandlerFunc provides basic JSON-RPC handling.
func HandlerFunc(w http.ResponseWriter, r *http.Request) {

	rs, batch, err := ParseRequest(r)
	if err != nil {
		SendResponse(w, []*Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		return
	}

	c := r.Context()
	resp := make([]*Response, len(rs))
	for i := range rs {
		resp[i] = invokeMethod(c, rs[i])
	}

	if err := SendResponse(w, resp, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func invokeMethod(c context.Context, r *Request) *Response {
	var h Handler
	res := NewResponse(r)
	h, res.Error = TakeMethod(r)
	if res.Error != nil {
		return res
	}
	res.Result, res.Error = h.ServeJSONRPC(WithRequestID(c, r.ID), r.Params)
	if res.Error != nil {
		res.Result = nil
	}
	return res
}
