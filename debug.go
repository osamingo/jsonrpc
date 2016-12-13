package jsonrpc

import (
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"runtime"

	"github.com/alecthomas/jsonschema"
)

// A MethodReference is a reference of JSON-RPC method.
type MethodReference struct {
	Name     string             `json:"name"`
	Function string             `json:"function"`
	Params   *jsonschema.Schema `json:"params,omitempty"`
	Result   *jsonschema.Schema `json:"result,omitempty"`
}

// DebugHandler views registered method list.
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeKey, contentTypeValue)
	ms := Methods()
	if len(ms) == 0 {
		w.Write([]byte("[]"))
		return
	}
	l := make([]MethodReference, 0, len(ms))
	for k, md := range ms {
		mr := MethodReference{
			Name:     k,
			Function: path.Base(runtime.FuncForPC(reflect.ValueOf(md.Func).Pointer()).Name()),
		}
		if md.Params != nil {
			mr.Params = jsonschema.Reflect(md.Params)
		}
		if md.Result != nil {
			mr.Result = jsonschema.Reflect(md.Result)
		}
		l = append(l, mr)
	}
	b, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
