package jsonrpc

import (
	"context"
	"testing"

	"github.com/intel-go/fastjson"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {

	c := context.Background()
	id := fastjson.RawMessage("1")
	c = WithRequestID(c, &id)
	var pick *fastjson.RawMessage
	require.NotPanics(t, func() {
		pick = RequestID(c)
	})
	require.Equal(t, &id, pick)
}

func TestMetadata(t *testing.T) {

	c := context.Background()
	md := Metadata{Params: Metadata{}}
	c = WithMetadata(c, md)
	var pick Metadata
	require.NotPanics(t, func() {
		pick = GetMetadata(c)
	})
	require.Equal(t, md, pick)
}

func TestMethodName(t *testing.T) {

	c := context.Background()
	c = WithMethodName(c, t.Name())
	var pick string
	require.NotPanics(t, func() {
		pick = MethodName(c)
	})
	require.Equal(t, t.Name(), pick)
}
