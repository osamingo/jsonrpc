package jsonrpc

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intel-go/fastjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {

	r, rerr := http.NewRequest("", "", bytes.NewReader(nil))
	require.NoError(t, rerr)

	_, _, err := ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("test")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("{}")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err := ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.False(t, batch)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("[")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("[test]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	_, _, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, rerr = http.NewRequest("", "", bytes.NewReader([]byte("[{}]")))
	require.NoError(t, rerr)

	r.Header.Set("Content-Type", "application/json")
	rs, batch, err = ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
	assert.True(t, batch)
}

func TestNewResponse(t *testing.T) {
	id := fastjson.RawMessage("test")
	r := NewResponse(&Request{
		Version: "2.0",
		ID:      &id,
	})
	assert.Equal(t, "2.0", r.Version)
	assert.Equal(t, "test", string(*r.ID))
}

func TestSendResponse(t *testing.T) {

	rec := httptest.NewRecorder()
	err := SendResponse(rec, []*Response{}, false)
	require.NoError(t, err)
	assert.Empty(t, rec.Body.String())

	id := fastjson.RawMessage([]byte(`"test"`))
	r := &Response{
		ID:      &id,
		Version: "2.0",
		Result: struct {
			Name string `json:"name"`
		}{
			Name: "john",
		},
	}

	rec = httptest.NewRecorder()
	err = SendResponse(rec, []*Response{r}, false)
	require.NoError(t, err)
	assert.Equal(t, `{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = SendResponse(rec, []*Response{r}, true)
	require.NoError(t, err)
	assert.Equal(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = SendResponse(rec, []*Response{r, r}, false)
	require.NoError(t, err)
	assert.Equal(t, `[{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"},{"jsonrpc":"2.0","result":{"name":"john"},"id":"test"}]
`, rec.Body.String())
}
