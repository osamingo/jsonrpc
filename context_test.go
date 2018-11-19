package jsonrpc

import (
	"context"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {

	c := context.Background()
	id := json.RawMessage("1")
	c = WithRequestID(c, &id)
	var pick *json.RawMessage
	require.NotPanics(t, func() {
		pick = RequestID(c)
	})
	require.Equal(t, &id, pick)
}
