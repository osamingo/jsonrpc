package jsonrpc

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {

	r, _ := http.NewRequest("", "", bytes.NewReader(nil))

	_, err := ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r.Header.Set("Content-Type", "application/json")

	_, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("")))
	r.Header.Set("Content-Type", "application/json")
	_, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeInvalidRequest, err.Code)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("test")))
	r.Header.Set("Content-Type", "application/json")
	_, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("{}")))
	r.Header.Set("Content-Type", "application/json")
	rs, err := ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("[")))
	r.Header.Set("Content-Type", "application/json")
	_, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("[test]")))
	r.Header.Set("Content-Type", "application/json")
	_, err = ParseRequest(r)
	require.IsType(t, &Error{}, err)
	assert.Equal(t, ErrorCodeParse, err.Code)

	r, _ = http.NewRequest("", "", bytes.NewReader([]byte("[{}]")))
	r.Header.Set("Content-Type", "application/json")
	rs, err = ParseRequest(r)
	require.Nil(t, err)
	require.NotEmpty(t, rs)
}

func TestNewResponse(t *testing.T) {
	r := NewResponse(Request{
		Version: "2.0",
		ID:      "test",
	})
	assert.Equal(t, "2.0", r.Version)
	assert.Equal(t, "test", r.ID)
}

func TestSendResponse(t *testing.T) {

	rec := httptest.NewRecorder()
	err := SendResponse(rec, []Response{})
	require.NoError(t, err)
	assert.Empty(t, rec.Body.String())

	r := Response{
		ID:      "test",
		Version: "2.0",
		Result: struct {
			Name string `json:"name"`
		}{
			Name: "john",
		},
	}

	rec = httptest.NewRecorder()
	err = SendResponse(rec, []Response{r})
	require.NoError(t, err)
	assert.Equal(t, `{"id":"test","jsonrpc":"2.0","result":{"name":"john"}}`, rec.Body.String())

	rec = httptest.NewRecorder()
	err = SendResponse(rec, []Response{r, r})
	require.NoError(t, err)
	assert.Equal(t, `[{"id":"test","jsonrpc":"2.0","result":{"name":"john"}},{"id":"test","jsonrpc":"2.0","result":{"name":"john"}}]`, rec.Body.String())
}
