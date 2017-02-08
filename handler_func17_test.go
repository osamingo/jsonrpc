// +build go1.7

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type handler struct {
	F func(c context.Context, params *json.RawMessage) (interface{}, *Error)
}

func (h handler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *Error) {
	return h.F(c, params)
}

func TestHandler17(t *testing.T) {

	PurgeMethods()

	rec := httptest.NewRecorder()
	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	HandlerFunc(rec, r)

	res := Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	HandlerFunc(rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	h := handler{}
	h.F = func(c context.Context, params *json.RawMessage) (interface{}, *Error) {
		return "hello", nil
	}
	require.NoError(t, RegisterMethod("hello", h, nil, nil))

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	HandlerFunc(rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.Nil(t, res.Error)
	assert.Equal(t, "hello", res.Result)
}
