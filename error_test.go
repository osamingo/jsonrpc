package jsonrpc_test

import (
	"testing"

	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Parallel()

	var err any = &jsonrpc.Error{}
	_, ok := err.(error)
	require.True(t, ok)
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	err := &jsonrpc.Error{
		Code:    jsonrpc.ErrorCode(100),
		Message: "test",
		Data: map[string]string{
			"test": "test",
		},
	}

	assert.Equal(t, "jsonrpc: code: 100, message: test, data: map[test:test]", err.Error())
}

func TestErrParse(t *testing.T) {
	t.Parallel()

	err := jsonrpc.ErrParse()
	require.Equal(t, jsonrpc.ErrorCodeParse, err.Code)
}

func TestErrInvalidRequest(t *testing.T) {
	t.Parallel()

	err := jsonrpc.ErrInvalidRequest()
	require.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)
}

func TestErrMethodNotFound(t *testing.T) {
	t.Parallel()

	err := jsonrpc.ErrMethodNotFound()
	require.Equal(t, jsonrpc.ErrorCodeMethodNotFound, err.Code)
}

func TestErrInvalidParams(t *testing.T) {
	t.Parallel()

	err := jsonrpc.ErrInvalidParams()
	require.Equal(t, jsonrpc.ErrorCodeInvalidParams, err.Code)
}

func TestErrInternal(t *testing.T) {
	t.Parallel()

	err := jsonrpc.ErrInternal()
	require.Equal(t, jsonrpc.ErrorCodeInternal, err.Code)
}
