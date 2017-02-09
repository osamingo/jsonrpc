package jsonrpc

import (
	"errors"
	"sync"
)

type (
	// A MethodRepository has JSON-RPC method functions.
	MethodRepository struct {
		m sync.RWMutex
		r map[string]Metadata
	}
	// Metadata has method meta data.
	Metadata struct {
		Handler Handler
		Params  interface{}
		Result  interface{}
	}
)

var mr = MethodRepository{
	m: sync.RWMutex{},
	r: map[string]Metadata{},
}

// TakeMethod takes jsonrpc.Func in MethodRepository.
func TakeMethod(r *Request) (Handler, *Error) {
	if r.Method == "" || r.Version != Version {
		return nil, ErrInvalidParams()
	}

	mr.m.RLock()
	md, ok := mr.r[r.Method]
	mr.m.RUnlock()
	if !ok {
		return nil, ErrMethodNotFound()
	}

	return md.Handler, nil
}

// RegisterMethod registers jsonrpc.Func to MethodRepository.
func RegisterMethod(method string, h Handler, params, result interface{}) error {
	if method == "" || h == nil {
		return errors.New("jsonrpc: method name and function should not be empty")
	}
	mr.m.Lock()
	mr.r[method] = Metadata{
		Handler: h,
		Params:  params,
		Result:  result,
	}
	mr.m.Unlock()
	return nil
}

// Methods returns registered methods.
func Methods() map[string]Metadata {
	mr.m.RLock()
	ml := make(map[string]Metadata, len(mr.r))
	for k, md := range mr.r {
		ml[k] = md
	}
	mr.m.RUnlock()
	return ml
}

// PurgeMethods purges all registered methods.
func PurgeMethods() {
	mr.m.Lock()
	mr.r = map[string]Metadata{}
	mr.m.Unlock()
}
