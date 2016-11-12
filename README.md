# jsonrpc

[![Travis branch](https://img.shields.io/travis/osamingo/jsonrpc/master.svg)](https://travis-ci.org/osamingo/jsonrpc)
[![codecov](https://codecov.io/gh/osamingo/jsonrpc/branch/master/graph/badge.svg)](https://codecov.io/gh/osamingo/jsonrpc)
[![Go Report Card](https://goreportcard.com/badge/osamingo/jsonrpc)](https://goreportcard.com/report/osamingo/jsonrpc)
[![codebeat badge](https://codebeat.co/badges/cbd0290d-200b-4693-80dc-296d9447c35b)](https://codebeat.co/projects/github-com-osamingo-jsonrpc)
[![GoDoc](https://godoc.org/github.com/osamingo/jsonrpc?status.svg)](https://godoc.org/github.com/osamingo/jsonrpc)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/osamingo/jsonrpc/master/LICENSE)

## About

- Simple implements ;)
- No `reflect` package.
- Support both packages `context` and `golang.org/x/net/context`.
- Compliance with [JSON-RPC 2.0](http://www.jsonrpc.org/specification).

## Install

```
$ go get -u github.com/osamingo/jsonrpc
```

## Usage

```go
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
```

## License

Released under the [MIT License](https://github.com/osamingo/jsonrpc/blob/master/LICENSE).
