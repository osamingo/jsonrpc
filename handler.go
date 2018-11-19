package jsonrpc

import (
	"context"
	"net/http"

	json "github.com/json-iterator/go"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c context.Context, params *json.RawMessage) (result interface{}, err *Error)
}

// ServeHTTP provides basic JSON-RPC handling.
func (mr *MethodRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rs, batch, err := ParseRequest(r)
	if err != nil {
		err := SendResponse(w, []*Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := make([]*Response, len(rs))
	for i := range rs {
		resp[i] = mr.InvokeMethod(r.Context(), rs[i])
	}

	if err := SendResponse(w, resp, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// InvokeMethod invokes JSON-RPC method.
func (mr *MethodRepository) InvokeMethod(c context.Context, r *Request) *Response {
	var h Handler
	res := NewResponse(r)
	h, res.Error = mr.TakeMethod(r)
	if res.Error != nil {
		return res
	}
	res.Result, res.Error = h.ServeJSONRPC(WithRequestID(c, r.ID), r.Params)
	if res.Error != nil {
		res.Result = nil
	}
	return res
}
