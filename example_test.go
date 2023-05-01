package jsonrpc_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/goccy/go-json"
	"github.com/osamingo/jsonrpc/v2"
)

type (
	EchoHandler struct{}
	EchoParams  struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}

	PositionalHandler struct{}
	PositionalParams  []int
	PositionalResult  struct {
		Message []int `json:"message"`
	}
)

func (h EchoHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var p EchoParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func (h PositionalHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {
	var p PositionalParams
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return PositionalResult{
		Message: p,
	}, nil
}

func ExampleMethodRepository_ServeHTTP() { //nolint: nosnakecase
	mr := jsonrpc.NewMethodRepository()

	if err := mr.RegisterMethod("Main.Echo", EchoHandler{}, EchoParams{}, EchoResult{}); err != nil {
		log.Println(err)

		return
	}

	if err := mr.RegisterMethod("Main.Positional", PositionalHandler{}, PositionalParams{}, PositionalResult{}); err != nil {
		log.Println(err)

		return
	}

	http.Handle("/jrpc", mr)
	http.HandleFunc("/jrpc/debug", mr.ServeDebug)

	srv := httptest.NewServer(http.DefaultServeMux)
	defer srv.Close()

	contextType := "application/json"
	echoVal := `{
		"jsonrpc": "2.0",
		"method": "Main.Echo",
		"params": {
			"name": "John Doe"
		},
		"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
	}`

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/jrpc", bytes.NewBufferString(echoVal))
	if err != nil {
		log.Println(err)

		return
	}

	req.Header.Add("Content-Type", contextType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Println(err)

		return
	}

	positionalVal := `{
		"jsonrpc": "2.0",
		"method": "Main.Positional",
		"params": [3, 1, 1, 3, 5, 3],
		"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
	}`

	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/jrpc", bytes.NewBufferString(positionalVal))
	if err != nil {
		log.Println(err)

		return
	}

	req.Header.Add("Content-Type", contextType)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Println(err)

		return
	}

	// Output:
	// {"jsonrpc":"2.0","result":{"message":"Hello, John Doe"},"id":"243a718a-2ebb-4e32-8cc8-210c39e8a14b"}
	// {"jsonrpc":"2.0","result":{"message":[3,1,1,3,5,3]},"id":"243a718a-2ebb-4e32-8cc8-210c39e8a14b"}
}
