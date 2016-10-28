// +build go1.7

package jsonrpc

import "net/http"

// Handler provides basic JSON-RPC handling.
func Handler(w http.ResponseWriter, r *http.Request) {

	rs, batch, err := ParseRequest(r)
	if err != nil {
		SendResponse(w, []Response{
			{
				Version: Version,
				Error:   err,
			},
		}, false)
		return
	}

	resp := make([]Response, len(rs))
	for i := range rs {
		var f Func
		res := NewResponse(rs[i])
		f, res.Error = TakeMethod(rs[i])
		if res.Error != nil {
			resp[i] = res
			continue
		}

		res.Result, res.Error = f(r.Context(), rs[i].Params)
		resp[i] = res
	}

	if err := SendResponse(w, resp, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
