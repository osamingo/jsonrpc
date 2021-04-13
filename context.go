package jsonrpc

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
)

type (
	requestIDKey  struct{}
	metadataIDKey struct{}
	methodNameKey struct{}
	requestKey    struct{}
	responseKey   struct{}
)

// RequestID takes request id from context.
func RequestID(c context.Context) *json.RawMessage {
	return c.Value(requestIDKey{}).(*json.RawMessage)
}

// WithRequestID adds request id to context.
func WithRequestID(c context.Context, id *json.RawMessage) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}

// GetMetadata takes jsonrpc metadata from context.
func GetMetadata(c context.Context) Metadata {
	return c.Value(metadataIDKey{}).(Metadata)
}

// WithMetadata adds jsonrpc metadata to context.
func WithMetadata(c context.Context, md Metadata) context.Context {
	return context.WithValue(c, metadataIDKey{}, md)
}

// MethodName takes method name from context.
func MethodName(c context.Context) string {
	return c.Value(methodNameKey{}).(string)
}

// WithMethodName adds method name to context.
func WithMethodName(c context.Context, name string) context.Context {
	return context.WithValue(c, methodNameKey{}, name)
}

// WithRequest adds request to context.
func WithRequest(c context.Context, r *http.Request) context.Context {
	return context.WithValue(c, requestKey{}, r)
}

// GetRequest takes request from context.
func GetRequest(c context.Context) *http.Request {
	v := c.Value(requestKey{})
	if r, ok := v.(*http.Request); ok {
		return r
	}

	return nil
}

// WithResponse adds response to context.
func WithResponse(c context.Context, r http.ResponseWriter) context.Context {
	return context.WithValue(c, responseKey{}, r)
}

// GetResponse takes response from context.
func GetResponse(c context.Context) http.ResponseWriter {
	v := c.Value(responseKey{})
	if r, ok := v.(http.ResponseWriter); ok {
		return r
	}

	return nil
}
