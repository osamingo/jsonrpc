package jsonrpc

import (
	"errors"
	"sync"
	"runtime"
	"reflect"
	"path"
)

// A MethodRepository has JSON-RPC method functions.
type MethodRepository struct {
	m sync.RWMutex
	r map[string]Func
}

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

// MethodList returns registered method list.
func MethodList() map[string]string {
	mr.m.RLock()
	ml := make(map[string]string, len(mr.r))
	for k, f := range mr.r {
		ml[k] = path.Base(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	}
	mr.m.RUnlock()
	return ml
}
