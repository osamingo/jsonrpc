package jsonrpc_test

import (
	"context"
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeMethod(t *testing.T) {
	t.Parallel()

	mr := jsonrpc.NewMethodRepository()

	r := &jsonrpc.Request{}
	_, err := mr.TakeMethod(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidParams, err.Code)

	r.Method = "test"
	_, err = mr.TakeMethod(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidParams, err.Code)

	r.Version = "2.0"
	_, err = mr.TakeMethod(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeMethodNotFound, err.Code)

	require.NoError(t, mr.RegisterMethod("test", SampleHandler(), nil, nil))

	f, err := mr.TakeMethod(r)
	require.Nil(t, err)
	assert.NotEmpty(t, f)
}

func TestRegisterMethod(t *testing.T) {
	t.Parallel()

	mr := jsonrpc.NewMethodRepository()

	err := mr.RegisterMethod("", nil, nil, nil)
	require.Error(t, err)

	err = mr.RegisterMethod("test", nil, nil, nil)
	require.Error(t, err)

	err = mr.RegisterMethod("test", SampleHandler(), nil, nil)
	require.NoError(t, err)
}

func TestMethods(t *testing.T) {
	t.Parallel()

	mr := jsonrpc.NewMethodRepository()

	err := mr.RegisterMethod("JsonRpc.Sample", SampleHandler(), nil, nil)
	require.NoError(t, err)

	ml := mr.Methods()
	require.NotEmpty(t, ml)
	assert.NotEmpty(t, ml["JsonRpc.Sample"].Handler)
}

func SampleHandler() *jsonrpc.HandlerFunc {
	h := jsonrpc.HandlerFunc(func(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
		return (any)(nil), nil
	})

	return &h
}
