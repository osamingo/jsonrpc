package jsonrpc

import (
	"encoding/json"
	"net/http"
)

// DebugHandler views registered method list.
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeKey, contentTypeValue)
	m := Methods()
	if len(m) == 0 {
		w.Write([]byte("[]"))
		return
	}
	l := make([]Metadata, 0, len(m))
	for k := range m {
		l = append(l, Metadata{
			Method:   k,
			Function: m[k],
		})
	}
	b, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
