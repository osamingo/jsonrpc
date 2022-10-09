package jsonrpc_test

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	err := jsonrpc.Unmarshal(nil, nil)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidParams, err.Code)

	src := json.RawMessage([]byte(`{"name":"john"}`))

	err = jsonrpc.Unmarshal(&src, nil)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidParams, err.Code)

	dst := struct {
		Name string `json:"name"`
	}{}

	err = jsonrpc.Unmarshal(&src, &dst)
	require.Nil(t, err)
	assert.Equal(t, "john", dst.Name)
}
