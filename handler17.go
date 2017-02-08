// +build go1.7

package jsonrpc

import (
	"context"
	"encoding/json"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c context.Context, params *json.RawMessage) (result interface{}, err *Error)
}
