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

	require.Equal(t, http.StatusNotFound, rec.Code)

	require.NoError(t, RegisterMethod("Debug.Sample", SampleHandler(), struct {
		Name string `json:"name"`
	}{}, struct {
		Message string `json:"message,omitrmpty"`
	}{}))

	rec = httptest.NewRecorder()
	r, err = http.NewRequest("", "", nil)
	require.NoError(t, err)

	DebugHandler(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.String())
}
