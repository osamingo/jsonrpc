package jsonrpc

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	mr := NewMethodRepository()

	rec := httptest.NewRecorder()
	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	mr.ServeHTTP(rec, r)

	res := Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)

	h1 := HandlerFunc(func(c context.Context, params *json.RawMessage) (interface{}, *Error) {
		return "hello", nil
	})
	require.NoError(t, mr.RegisterMethod("hello", h1, nil, nil))
	h2 := HandlerFunc(func(c context.Context, params *json.RawMessage) (interface{}, *Error) {
		return nil, ErrInternal()
	})
	require.NoError(t, mr.RegisterMethod("bye", h2, nil, nil))

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"hello","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.Nil(t, res.Error)
	assert.Equal(t, "hello", res.Result)

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":"test","method":"bye","params":{}}`)))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	mr.ServeHTTP(rec, r)
	res = Response{}
	err = json.NewDecoder(rec.Body).Decode(&res)
	require.NoError(t, err)
	assert.NotNil(t, res.Error)
}

func TestInvokeMethodMiddlewares(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := context.Background()
	id := json.RawMessage("test")
	r := &Request{
		Version: "2.0",
		Method:  "test",
		ID:      &id,
	}

	mr := NewMethodRepository()
	err := mr.RegisterMethod("test", HandlerFunc(func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
		v := c.Value("key1")
		require.NotNil(t, v)
		v = c.Value("key2")
		require.NotNil(t, v)
		return "value3", nil
	}), nil, nil, func(next HandlerFunc) HandlerFunc {
		return func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
			c = context.WithValue(c, "key1", "value1")
			return next(c, params)
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
			v := c.Value("key1")
			require.NotNil(t, v)
			c = context.WithValue(c, "key2", "value2")
			return next(c, params)
		}
	})
	require.NoError(t, err)

	resp := mr.InvokeMethod(ctx, r, req, rec)
	require.Nil(t, resp.Error)
	require.NotNil(t, resp.Result)
}
