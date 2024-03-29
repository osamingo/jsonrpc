package jsonrpc

import "fmt"

const (
	// ErrorCodeParse is parse error code.
	ErrorCodeParse ErrorCode = -32700
	// ErrorCodeInvalidRequest is invalid request error code.
	ErrorCodeInvalidRequest ErrorCode = -32600
	// ErrorCodeMethodNotFound is method not found error code.
	ErrorCodeMethodNotFound ErrorCode = -32601
	// ErrorCodeInvalidParams is invalid params error code.
	ErrorCodeInvalidParams ErrorCode = -32602
	// ErrorCodeInternal is internal error code.
	ErrorCodeInternal ErrorCode = -32603
)

type (
	// A ErrorCode by JSON-RPC 2.0.
	ErrorCode int

	// An Error is a wrapper for a JSON interface value.
	Error struct {
		Code    ErrorCode `json:"code"`
		Message string    `json:"message"`
		Data    any       `json:"data,omitempty"`
	}
)

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc: code: %d, message: %s, data: %+v", e.Code, e.Message, e.Data)
}

// ErrParse returns parse error.
func ErrParse() *Error {
	return &Error{
		Code:    ErrorCodeParse,
		Message: "Parse error",
	}
}

// ErrInvalidRequest returns invalid request error.
func ErrInvalidRequest() *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: "Invalid Request",
	}
}

// ErrMethodNotFound returns method not found error.
func ErrMethodNotFound() *Error {
	return &Error{
		Code:    ErrorCodeMethodNotFound,
		Message: "Method not found",
	}
}

// ErrInvalidParams returns invalid params error.
func ErrInvalidParams() *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: "Invalid params",
	}
}

// ErrInternal returns internal error.
func ErrInternal() *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "Internal error",
	}
}
