package jsonrpc_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	t.Parallel()

	r, rerr := http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader(nil))
	require.NoError(t, rerr)

	_, _, err := jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("test")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("{}")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err := jsonrpc.ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.False(t, batch)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("[")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("[test]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = jsonrpc.ParseRequest(r)
	require.IsType(t, &jsonrpc.Error{}, err)
	assert.Equal(t, jsonrpc.ErrorCodeParse, err.Code)

	r, rerr = http.NewRequestWithContext(t.Context(), "", "", bytes.NewReader([]byte("[{}]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err = jsonrpc.ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.True(t, batch)
}

func TestParseRequestUnknownContentLength(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		body      string
		batch     bool
		errorCode jsonrpc.ErrorCode
	}{
		{name: "single", body: `{"jsonrpc":"2.0","method":"hello","id":1}`},
		{name: "batch", body: `[{"jsonrpc":"2.0","method":"hello","id":1}]`, batch: true},
		{name: "malformed", body: "[", errorCode: jsonrpc.ErrorCodeParse},
		{name: "empty", errorCode: jsonrpc.ErrorCodeInvalidRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/jrpc", io.NopCloser(strings.NewReader(tc.body)))
			r.Header.Set("Content-Type", "application/json")
			require.Equal(t, int64(-1), r.ContentLength)

			var (
				requests []*jsonrpc.Request
				batch    bool
				err      *jsonrpc.Error
			)
			require.NotPanics(t, func() {
				requests, batch, err = jsonrpc.ParseRequest(r)
			})
			if tc.errorCode != 0 {
				require.NotNil(t, err)
				assert.Equal(t, tc.errorCode, err.Code)

				return
			}

			require.Nil(t, err)
			require.Len(t, requests, 1)
			assert.Equal(t, tc.batch, batch)
			assert.Equal(t, "hello", requests[0].Method)
		})
	}
}

func TestHandlerContentLength(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		body     string
		response string
	}{
		{
			name:     "single",
			body:     `{"jsonrpc":"2.0","method":"hello","id":1}`,
			response: `{"jsonrpc":"2.0","result":"hello","id":1}`,
		},
		{
			name:     "batch",
			body:     `[{"jsonrpc":"2.0","method":"hello","id":1}]`,
			response: `[{"jsonrpc":"2.0","result":"hello","id":1}]`,
		},
	} {
		for _, chunked := range []bool{false, true} {
			name := tc.name + "/known"
			if chunked {
				name = tc.name + "/chunked"
			}

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				mr := jsonrpc.NewMethodRepository()
				h := jsonrpc.HandlerFunc(func(context.Context, *json.RawMessage) (any, *jsonrpc.Error) {
					return "hello", nil
				})
				require.NoError(t, mr.RegisterMethod("hello", h, nil, nil))

				lengths := make(chan int64, 1)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					lengths <- r.ContentLength
					mr.ServeHTTP(w, r)
				}))
				t.Cleanup(server.Close)

				var body io.Reader = strings.NewReader(tc.body)
				if chunked {
					body = io.NopCloser(body)
				}
				r, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL, body)
				require.NoError(t, err)
				r.Header.Set("Content-Type", "application/json")

				res, err := server.Client().Do(r)
				require.NoError(t, err)
				t.Cleanup(func() { assert.NoError(t, res.Body.Close()) })
				got, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.JSONEq(t, tc.response, string(got))

				wantLength := int64(len(tc.body))
				if chunked {
					wantLength = -1
				}
				assert.Equal(t, wantLength, <-lengths)
			})
		}
	}
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
	assert.JSONEq(t, `{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = jsonrpc.SendResponse(rec, []*jsonrpc.Response{r}, true)
	require.NoError(t, err)
	assert.JSONEq(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = jsonrpc.SendResponse(rec, []*jsonrpc.Response{r, r}, false)
	require.NoError(t, err)
	assert.JSONEq(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"},{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())
}
