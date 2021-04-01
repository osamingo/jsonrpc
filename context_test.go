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

func TestMetadata(t *testing.T) {
	c := context.Background()
	md := Metadata{Params: Metadata{}}
	c = WithMetadata(c, md)
	var pick Metadata
	require.NotPanics(t, func() {
		pick = GetMetadata(c)
	})
	require.Equal(t, md, pick)
}

func TestMethodName(t *testing.T) {
	c := context.Background()
	c = WithMethodName(c, t.Name())
	var pick string
	require.NotPanics(t, func() {
		pick = MethodName(c)
	})
	require.Equal(t, t.Name(), pick)
}

func TestRequest(t *testing.T) {
	assert.NotPanics(t, func() {
		r := GetRequest(context.Background())
		assert.Nil(t, r)
	})
	c := context.Background()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	c = WithRequest(c, r)
	assert.Equal(t, r, GetRequest(c))
}

func TestResponse(t *testing.T) {
	assert.NotPanics(t, func() {
		r := GetResponse(context.Background())
		assert.Nil(t, r)
	})
	c := context.Background()
	r := httptest.NewRecorder()
	c = WithResponse(c, r)
	assert.Equal(t, r, GetResponse(c))
}
