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

func TestParseRequest(t *testing.T) {
	t.Parallel()

	r, rerr := http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader(nil))
	require.NoError(t, rerr)

	_, _, err := jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("test")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("{}")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err := jsonrpc.ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.False(t, batch)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("[")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("[test]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(context.Background(), "", "", bytes.NewReader([]byte("[{}]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err = jsonrpc.ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.True(t, batch)
}

func TestNewResponse(t *testing.T) {
	t.Parallel()

	id := json.RawMessage("test")
	r := jsonrpc.NewResponse(&jsonrpc.Request{
		Version: "2.0",
		ID:      &id,
	})
	assert.Equal(t, "2.0", r.Version)
	assert.Equal(t, "test", string(*r.ID))
}

func TestSendResponse(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	err := jsonrpc.SendResponse(rec, []*jsonrpc.Response{}, false)
	require.NoError(t, err)
	assert.Empty(t, rec.Body.String())

	id := json.RawMessage([]byte(`"test"`))
	r := &jsonrpc.Response{
		ID:      &id,
		Version: "2.0",
		Result: struct {
			Name string `json:"name"`
		}{
			Name: "john",
		},
	}

	rec = httptest.NewRecorder()
	err = jsonrpc.SendResponse(rec, []*jsonrpc.Response{r}, false)
	require.NoError(t, err)
	assert.Equal(t, `{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = jsonrpc.SendResponse(rec, []*jsonrpc.Response{r}, true)
	require.NoError(t, err)
	assert.Equal(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = jsonrpc.SendResponse(rec, []*jsonrpc.Response{r, r}, false)
	require.NoError(t, err)
	assert.Equal(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"},{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())
}
