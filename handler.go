package jsonrpc

import (
	"encoding/json"

	"golang.org/x/net/context"
)

// Handler links a method of JSON-RPC request.
type Handler interface {
	ServeJSONRPC(c context.Context, params *json.RawMessage) (result interface{}, err *Error)
}
