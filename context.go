package jsonrpc

import (
	"context"

	"github.com/intel-go/fastjson"
)

type requestIDKey struct{}
type metadataIDKey struct{}

// RequestID takes request id from context.
func RequestID(c context.Context) *fastjson.RawMessage {
	return c.Value(requestIDKey{}).(*fastjson.RawMessage)
}

// WithRequestID adds request id to context.
func WithRequestID(c context.Context, id *fastjson.RawMessage) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}

// RequestID takes request id from context.
func GetMetadata(c context.Context) Metadata {
	return c.Value(metadataIDKey{}).(Metadata)
}

// WithRequestID adds request id to context.
func WithMetadata(c context.Context, md Metadata) context.Context {
	return context.WithValue(c, metadataIDKey{}, md)
}
