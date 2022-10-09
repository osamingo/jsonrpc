package jsonrpc_test

import (
	"context"
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	c := context.Background()
	id := json.RawMessage("1")
	c = jsonrpc.WithRequestID(c, &id)
	var pick *json.RawMessage
	require.NotPanics(t, func() {
		pick = jsonrpc.RequestID(c)
	})
	require.Equal(t, &id, pick)
}

func TestMetadata(t *testing.T) {
	t.Parallel()

	c := context.Background()
	md := jsonrpc.Metadata{Params: jsonrpc.Metadata{}}
	c = jsonrpc.WithMetadata(c, md)
	var pick jsonrpc.Metadata
	require.NotPanics(t, func() {
		pick = jsonrpc.GetMetadata(c)
	})
	require.Equal(t, md, pick)
}

func TestMethodName(t *testing.T) {
	t.Parallel()

	c := context.Background()
	c = jsonrpc.WithMethodName(c, t.Name())
	var pick string
	require.NotPanics(t, func() {
		pick = jsonrpc.MethodName(c)
	})
	require.Equal(t, t.Name(), pick)
}
