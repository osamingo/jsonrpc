package jsonrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseMiddleware(t *testing.T) {
	mr := NewMethodRepository()
	assert.Len(t, mr.middlewares, 0)
	mr.UseMiddleware(nil)
	assert.Len(t, mr.middlewares, 1)
	mr.UseMiddleware(nil, nil)
	assert.Len(t, mr.middlewares, 3)
}

func TestMiddlewareOrder(t *testing.T) {
	key := "key"
	mw := func(want int) MiddlewareFunc {
		return func(next HandlerFunc) HandlerFunc {
			return func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
				got := c.Value(key)
				v, ok := got.(int)
				if assert.True(t, ok) {
					assert.Equal(t, want, v)
				}

				c = context.WithValue(c, key, v+1)
				return next(c, params)
			}
		}
	}

	mr := NewMethodRepository()
	mr.UseMiddleware(mw(3), mw(2))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := context.Background()
	id := json.RawMessage("test")
	r := &Request{
		Version: "2.0",
		Method:  "test",
		ID:      &id,
	}

	err := mr.RegisterMethod("test", HandlerFunc(func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
		got := c.Value(key)
		v, ok := got.(int)
		require.True(t, ok)
		assert.Equal(t, 4, v)
		return nil, nil
	}), nil, nil, mw(1), mw(0))
	require.NoError(t, err)

	resp := mr.InvokeMethod(ctx, r, req, rec)
	require.Nil(t, resp.Error)
}
