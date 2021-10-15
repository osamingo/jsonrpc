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

func TestUseMiddleware(t *testing.T) {
	mr := NewMethodRepository()
	assert.Len(t, mr.middlewares, 0)
	mr.UseMiddleware(nil)
	assert.Len(t, mr.middlewares, 1)
	mr.UseMiddleware(nil, nil)
	assert.Len(t, mr.middlewares, 3)
}

func TestMiddlewareOrder(t *testing.T) {
	buf := bytes.Buffer{}
	mw := func(s string) MiddlewareFunc {
		return func(next HandlerFunc) HandlerFunc {
			return func(c context.Context, params *json.RawMessage) (result interface{}, err *Error) {
				buf.WriteString(s)
				return next(c, params)
			}
		}
	}

	mr := NewMethodRepository()
	mr.UseMiddleware(mw("-1"), mw("1"))
	mr.UseMiddleware(mw("2"))

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
		return nil, nil
	}),
		nil,
		nil,
		mw("3"),
		mw("4"),
		mw("5"),
	)
	require.NoError(t, err)

	resp := mr.InvokeMethod(ctx, r, req, rec)
	require.Nil(t, resp.Error)
	assert.Equal(t, "-112345", buf.String())
}
