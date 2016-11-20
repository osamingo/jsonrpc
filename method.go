package jsonrpc

import (
	"errors"
	"path"
	"reflect"
	"runtime"
	"sync"
)

type (
	// A MethodRepository has JSON-RPC method functions.
	MethodRepository struct {
		m sync.RWMutex
		r map[string]Func
	}
	// Metadata has method name and function name.
	Metadata struct {
		Method   string `json:"method"`
		Function string `json:"function"`
	}
)

var mr = MethodRepository{
	m: sync.RWMutex{},
	r: map[string]Func{},
}

// TakeMethod takes jsonrpc.Func in MethodRepository.
func TakeMethod(r Request) (Func, *Error) {
	if r.Method == "" || r.Version != Version {
		return nil, ErrInvalidParams()
	}

	mr.m.RLock()
	f, ok := mr.r[r.Method]
	mr.m.RUnlock()
	if !ok {
		return nil, ErrMethodNotFound()
	}

	return f, nil
}

// RegisterMethod registers jsonrpc.Func to MethodRepository.
func RegisterMethod(method string, f Func) error {
	if method == "" || f == nil {
		return errors.New("jsonrpc: method name and function should not be empty")
	}
	mr.m.Lock()
	mr.r[method] = f
	mr.m.Unlock()
	return nil
}

// Methods returns registered methods.
func Methods() map[string]string {
	mr.m.RLock()
	ml := make(map[string]string, len(mr.r))
	for k, f := range mr.r {
		ml[k] = path.Base(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	}
	mr.m.RUnlock()
	return ml
}

// PurgeMethods purges all registered methods.
func PurgeMethods() {
	mr.m.Lock()
	mr.r = map[string]Func{}
	mr.m.Unlock()
}
