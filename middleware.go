package jsonrpc

// UseMiddleware adds middlewares
func (mr *MethodRepository) UseMiddleware(middlewares ...MiddlewareFunc) {
	mr.middlewares = append(mr.middlewares, middlewares...)
}

// applyMiddleware applies middlewares to Handler
func applyMiddleware(h Handler, middleware ...MiddlewareFunc) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h.ServeJSONRPC)
	}
	return h
}
