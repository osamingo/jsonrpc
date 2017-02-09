// +build go1.7

package jsonrpc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeMethod(t *testing.T) {

	r := &Request{}
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

	require.NoError(t, RegisterMethod("test", SampleHandler(), nil, nil))

	f, err := TakeMethod(r)
	require.Nil(t, err)
	assert.NotEmpty(t, f)
}

func TestRegisterMethod(t *testing.T) {

	err := RegisterMethod("", nil, nil, nil)
	require.Error(t, err)

	err = RegisterMethod("test", nil, nil, nil)
	require.Error(t, err)

	err = RegisterMethod("test", SampleHandler(), nil, nil)
	require.NoError(t, err)
}

func TestMethods(t *testing.T) {

	err := RegisterMethod("JsonRpc.Sample", SampleHandler(), nil, nil)
	require.NoError(t, err)

	ml := Methods()
	require.NotEmpty(t, ml)
	assert.NotEmpty(t, ml["JsonRpc.Sample"].Handler)
}

func SampleHandler() Handler {
	h := handler{}
	h.F = func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
		return nil, nil
	}
	return h
}
