package jsonrpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osamingo/jsonrpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebugHandler(t *testing.T) {
	t.Parallel()

	mr := jsonrpc.NewMethodRepository()

	rec := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(t.Context(), "", "", nil)
	require.NoError(t, err)

	mr.ServeDebug(rec, r)

	require.Equal(t, http.StatusNotFound, rec.Code)

	require.NoError(t, mr.RegisterMethod("Debug.Sample", SampleHandler(), struct {
		Name string `json:"name"`
	}{}, struct {
		Message string `json:"message,omitempty"`
	}{}))

	rec = httptest.NewRecorder()
	r, err = http.NewRequestWithContext(t.Context(), "", "", nil)
	require.NoError(t, err)

	mr.ServeDebug(rec, r)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.String())
}
