package jsonrpc

import (
	"context"
	"testing"

	"github.com/intel-go/fastjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeMethodMetadata(t *testing.T) {

	mr := NewMethodRepository()

	r := &Request{}
	_, err := mr.TakeMethodMetadata(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidParams, err.Code)

	r.Method = "test"
	_, err = mr.TakeMethodMetadata(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidParams, err.Code)

	r.Version = "2.0"
	_, err = mr.TakeMethodMetadata(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeMethodNotFound, err.Code)

	require.NoError(t, mr.RegisterMethod("test", SampleHandler(), nil, nil))

	md, err := mr.TakeMethodMetadata(r)
	require.Nil(t, err)
	assert.NotEmpty(t, md)
}

func TestRegisterMethod(t *testing.T) {

	mr := NewMethodRepository()

	err := mr.RegisterMethod("", nil, nil, nil)
	require.Error(t, err)

	err = mr.RegisterMethod("test", nil, nil, nil)
	require.Error(t, err)

	err = mr.RegisterMethod("test", SampleHandler(), nil, nil)
	require.NoError(t, err)
}

func TestMethods(t *testing.T) {

	mr := NewMethodRepository()

	err := mr.RegisterMethod("JsonRpc.Sample", SampleHandler(), nil, nil)
	require.NoError(t, err)

	ml := mr.Methods()
	require.NotEmpty(t, ml)
	assert.NotEmpty(t, ml["JsonRpc.Sample"].Handler)
}

func SampleHandler() Handler {
	h := handler{}
	h.F = func(c context.Context, params *fastjson.RawMessage) (result interface{}, err *Error) {
		return nil, nil
	}
	return &h
}

type ExampleParam struct {
}

func TestMetadataInContext(t *testing.T) {

	mr := NewMethodRepository()

	err := mr.RegisterMethod(t.Name(), SampleHandler(), ExampleParam{}, nil)
	require.NoError(t, err)

	r := &Request{Method: t.Name(), Version: "2.0"}
	md, err := mr.TakeMethodMetadata(r)
	require.Nil(t, err)
	assert.Equal(t, ExampleParam{}, md.Params)

}
