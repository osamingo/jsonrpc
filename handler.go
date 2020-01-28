package jsonrpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/intel-go/fastjson"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c context.Context, params *fastjson.RawMessage) (result interface{}, err *Error)
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
			fmt.Fprint(w, "Failed to encode error objects")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := make([]*Response, len(rs))
	for i := range rs {
		resp[i] = mr.InvokeMethod(r.Context(), rs[i])
	}

	if err := SendResponse(w, resp, batch); err != nil {
		fmt.Fprint(w, "Failed to encode result objects")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// InvokeMethod invokes JSON-RPC method.
func (mr *MethodRepository) InvokeMethod(c context.Context, r *Request) *Response {
	var md Metadata
	res := NewResponse(r)
	md, res.Error = mr.TakeMethodMetadata(r)
	if res.Error != nil {
		return res
	}

	wrappedContext := WithMetadata(WithRequestID(c, r.ID), md)
	res.Result, res.Error = md.Handler.ServeJSONRPC(wrappedContext, r.Params)
	if res.Error != nil {
		res.Result = nil
	}
	return res
}
