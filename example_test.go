//+build !go1.6

package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"golang.org/x/net/context"
)

type (
	EchoHandler struct{}
	EchoParams  struct {
		Name string `json:"name"`
	}
	EchoResult struct {
		Message string `json:"message"`
	}
)

var _ (Handler) = (*EchoHandler)(nil)

func (h *EchoHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *Error) {

	var p EchoParams
	if err := Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return EchoResult{
		Message: "Hello, " + p.Name,
	}, nil
}

func ExampleEchoHandler_ServeJSONRPC() {

	if err := RegisterMethod("Main.Echo", &EchoHandler{}, EchoParams{}, EchoResult{}); err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/jrpc", func(w http.ResponseWriter, r *http.Request) {
		HandlerFunc(r.Context(), w, r)
	})
	http.HandleFunc("/jrpc/debug", DebugHandlerFunc)

	srv := httptest.NewServer(http.DefaultServeMux)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/jrpc", "application/json", bytes.NewBufferString(`{
	  "jsonrpc": "2.0",
      "method": "Main.Echo",
      "params": {
        "name": "John Doe"
      },
      "id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
    }`))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// {"id":"243a718a-2ebb-4e32-8cc8-210c39e8a14b","jsonrpc":"2.0","result":{"message":"Hello, John Doe"}}
}
