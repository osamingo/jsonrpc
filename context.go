// +build !go1.7

package jsonrpc

import (
	"encoding/json"

	"golang.org/x/net/context"
)

type requestIDKey struct{}

// RequestID takes request id from context.
func RequestID(c context.Context) *json.RawMessage {
	return c.Value(requestIDKey{}).(*json.RawMessage)
}

// WithRequestID adds request id to context.
func WithRequestID(c context.Context, id *json.RawMessage) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}
