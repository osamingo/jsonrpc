// +build go1.7

package jsonrpc

import (
	"context"
	"encoding/json"
)

type requestIDKey struct {}

func RequestID(c context.Context) *json.RawMessage {
	return c.Value(requestIDKey{}).(*json.RawMessage)
}

func WithRequestID(c context.Context, id *json.RawMessage) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}
