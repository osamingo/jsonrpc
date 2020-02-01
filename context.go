package jsonrpc

import (
	"context"

	"github.com/intel-go/fastjson"
)

type requestIDKey struct{}
type methodNameKey struct{}
type metadataIDKey struct{}

// RequestID takes request id from context.
func RequestID(c context.Context) *fastjson.RawMessage {
	return c.Value(requestIDKey{}).(*fastjson.RawMessage)
}

// WithRequestID adds request id to context.
func WithRequestID(c context.Context, id *fastjson.RawMessage) context.Context {
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
