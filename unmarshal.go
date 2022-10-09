package jsonrpc

import "github.com/goccy/go-json"

// Unmarshal decodes JSON-RPC params.
func Unmarshal(params *json.RawMessage, dst any) *Error {
	if params == nil {
		return ErrInvalidParams()
	}
	if err := json.Unmarshal(*params, dst); err != nil {
		return ErrInvalidParams()
	}

	return nil
}
