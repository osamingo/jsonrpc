package jsonrpc_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	mr := jsonrpc.NewMethodRepository()

	rec := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(t.Context(), "", "", nil)
	require.NoError(t, err)

	mr.ServeHTTP(rec, r)

	res := jsonrpc.Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	rec = httptest.NewRecorder()
	r, err = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = jsonrpc.Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	h1 := jsonrpc.HandlerFunc(func(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
		return "hello", nil
	})
	require.NoError(t, mr.RegisterMethod("hello", h1, nil, nil))
	h2 := jsonrpc.HandlerFunc(func(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
		return nil, jsonrpc.ErrInternal()
	})
	require.NoError(t, mr.RegisterMethod("bye", h2, nil, nil))

	rec = httptest.NewRecorder()
	r, err = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = jsonrpc.Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.Nil(t, res.Error)
	assert.Equal(t, "hello", res.Result)

	rec = httptest.NewRecorder()
	r, err = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"bye","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = jsonrpc.Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)
}
