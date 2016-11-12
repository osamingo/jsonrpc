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
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func init() {
	jsonrpc.RegisterMethod("Echo", Echo)
}

func main() {
	http.HandleFunc("/v1/jrpc", jsonrpc.Handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
