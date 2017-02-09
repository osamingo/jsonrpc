package jsonrpc

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/alecthomas/jsonschema"
)

// A MethodReference is a reference of JSON-RPC method.
type MethodReference struct {
	Name    string             `json:"name"`
	Handler string             `json:"handler"`
	Params  *jsonschema.Schema `json:"params,omitempty"`
	Result  *jsonschema.Schema `json:"result,omitempty"`
}

// DebugHandlerFunc views registered method list.
func DebugHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ms := Methods()
	if len(ms) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	l := make([]*MethodReference, 0, len(ms))
	for k, md := range ms {
		mr := &MethodReference{
			Name: k,
		}
		tv := reflect.TypeOf(md.Handler)
		if tv.Kind() == reflect.Ptr {
			tv = tv.Elem()
		}
		mr.Handler = tv.Name()
		if md.Params != nil {
			mr.Params = jsonschema.Reflect(md.Params)
		}
		if md.Result != nil {
			mr.Result = jsonschema.Reflect(md.Result)
		}
		l = append(l, mr)
	}
	w.Header().Set(contentTypeKey, contentTypeValue)
	if err := json.NewEncoder(w).Encode(l); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
