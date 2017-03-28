// +build go1.7

package jsonrpc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {

	c := context.Background()
	id := json.RawMessage("1")
	c = WithRequestID(c, &id)
	var pick *json.RawMessage
	require.NotPanics(t, func(){
		pick = RequestID(c)
	})
	require.Equal(t, &id, pick)
}
