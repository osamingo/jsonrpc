package jsonrpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c context.Context, params *json.RawMessage) (result interface{}, err *Error)
}

// HandlerFunc type is an adapter to allow the use of
// ordinary functions as JSONRPC handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// jsonrpc.Handler that calls f.
type HandlerFunc func(c context.Context, params *json.RawMessage) (result interface{}, err *Error)

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// ServeJSONRPC calls f(w, r).
func (f HandlerFunc) ServeJSONRPC(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
	return f(c, params)
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
		resp[i] = mr.InvokeMethod(r.Context(), rs[i], r, w)
	}

	if err := SendResponse(w, resp, batch); err != nil {
		fmt.Fprint(w, "Failed to encode result objects")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func applyMiddleware(h Handler, middleware ...MiddlewareFunc) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h.ServeJSONRPC)
	}
	return h
}

// InvokeMethod invokes JSON-RPC method.
func (mr *MethodRepository) InvokeMethod(c context.Context, r *Request, req *http.Request, w http.ResponseWriter) *Response {
	var md Metadata
	res := NewResponse(r)
	md, res.Error = mr.TakeMethodMetadata(r)
	if res.Error != nil {
		return res
	}

	wrappedContext := WithRequestID(c, r.ID)
	wrappedContext = WithMethodName(wrappedContext, r.Method)
	wrappedContext = WithMetadata(wrappedContext, md)
	wrappedContext = WithRequest(wrappedContext, req)
	wrappedContext = WithResponse(wrappedContext, w)
	handler := applyMiddleware(md.Handler, md.Middlewares...)
	res.Result, res.Error = handler.ServeJSONRPC(wrappedContext, r.Params)
	if res.Error != nil {
		res.Result = nil
	}
	return res
}
