package jsonrpc

import (
	"errors"
	"sync"
)

// A MethodRepository is management JSON-RPC method functions.
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
		return errors.New("jsonrpc: method and function should not be empty")
	}
	mr.m.Lock()
	mr.r[method] = f
	mr.m.Unlock()
	return nil
}
