package jsonrpc_test

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	id := json.RawMessage("1")
	c := jsonrpc.WithRequestID(t.Context(), &id)
	var pick *json.RawMessage
	require.NotPanics(t, func() {
		pick = jsonrpc.RequestID(c)
	})
	require.Equal(t, &id, pick)
}

func TestMetadata(t *testing.T) {
	t.Parallel()

	md := jsonrpc.Metadata{Params: jsonrpc.Metadata{}}
	c := jsonrpc.WithMetadata(t.Context(), md)
	var pick jsonrpc.Metadata
	require.NotPanics(t, func() {
		pick = jsonrpc.GetMetadata(c)
	})
	require.Equal(t, md, pick)
}

func TestMethodName(t *testing.T) {
	t.Parallel()

	c := jsonrpc.WithMethodName(t.Context(), t.Name())
	var pick string
	require.NotPanics(t, func() {
		pick = jsonrpc.MethodName(c)
	})
	require.Equal(t, t.Name(), pick)
}
