package jsonrpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestTakeMethod(t *testing.T) {

	r := Request{}
	_, err := TakeMethod(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidParams, err.Code)

	r.Method = "test"
	_, err = TakeMethod(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidParams, err.Code)

	r.Version = "2.0"
	_, err = TakeMethod(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeMethodNotFound, err.Code)

	require.NoError(t, RegisterMethod("test", func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
		return nil, nil
	}))

	f, err := TakeMethod(r)
	require.Nil(t, err)
	assert.NotEmpty(t, f)
}

func TestRegisterMethod(t *testing.T) {

	err := RegisterMethod("", nil)
	require.Error(t, err)

	err = RegisterMethod("test", nil)
	require.Error(t, err)

	err = RegisterMethod("test", func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
		return nil, nil
	})
	require.NoError(t, err)
}
