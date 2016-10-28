package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/osamingo/jsonrpc"
)

type (
	EchoParams struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}
)

func Echo(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p EchoParams
	if err := json.Unmarshal(*params, &p); err != nil {
		return nil, jsonrpc.ErrInvalidParams()
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func JSONRPC(w http.ResponseWriter, r *http.Request) {

	rs, err := jsonrpc.ParseRequest(r)
	if err != nil {
		jsonrpc.SendResponse(w, []jsonrpc.Response{
			{
				Version: jsonrpc.Version,
				Error:   err,
			},
		})
		return
	}

	resp := make([]jsonrpc.Response, 0, len(rs))
	for i := range rs {
		var f jsonrpc.Func
		res := jsonrpc.NewResponse(rs[i])
		f, res.Error = jsonrpc.TakeMethod(rs[i])
		if res.Error != nil {
			resp = append(resp, res)
			continue
		}

		res.Result, res.Error = f(r.Context(), rs[i].Params)
		resp = append(resp, res)
	}

	if err := jsonrpc.SendResponse(w, resp); err != nil {
		log.Println(err)
	}
}

func init() {
	jsonrpc.RegisterMethod("Echo", Echo)
}

func main() {
	http.HandleFunc("/_jr", JSONRPC)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
