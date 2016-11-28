// +build go1.7

package jsonrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebugHandler(t *testing.T) {

	PurgeMethods()

	rec := httptest.NewRecorder()
	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	DebugHandler(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]", rec.Body.String())

	require.NoError(t, RegisterMethod("Debug.Sample", SampleFunc, struct {
		Name string `json:"name"`
	}{}, struct {
		Message string `json:"message,omitrmpty"`
	}{}))

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", nil)
	require.NoError(t, err)

	DebugHandler(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `[{"name":"Debug.Sample","function":"jsonrpc.SampleFunc","params":{"$ref":"#/definitions/","definitions":{"":{"type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false,"required":["name"]}}},"result":{"$ref":"#/definitions/","definitions":{"":{"type":"object","properties":{"message":{"type":"string"}},"additionalProperties":false,"required":["message"]}}}}]`, rec.Body.String())
}
