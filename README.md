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
- Support GAE/Go Standard Environment.
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
	jsonrpc.RegisterMethod("Echo", Echo, EchoParams{}, EchoResult{})
}

func main() {
	http.HandleFunc("/v1/jrpc", jsonrpc.Handler)
	http.HandleFunc("/v1/jrpc/debug", jsonrpc.DebugHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
```

### Result

#### Invoke the Echo method

```
POST /v1/jrpc HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 82
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/0.9.6

{
  "jsonrpc": "2.0",
  "method": "Echo",
  "params": {
    "name": "John Doe"
  },
  "id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
}

HTTP/1.1 200 OK
Content-Length: 68
Content-Type: application/json
Date: Mon, 28 Nov 2016 13:48:13 GMT

{
  "jsonrpc": "2.0",
  "result": {
    "message": "Hello, John Doe"
  },
  "id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"
}
```

#### Access to debug handler

```
GET /v1/jrpc/debug HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/0.9.6



HTTP/1.1 200 OK
Content-Length: 408
Content-Type: application/json
Date: Mon, 28 Nov 2016 13:56:24 GMT

[
  {
    "name": "Echo",
    "function": "main.Echo",
    "params": {
      "$ref": "#/definitions/EchoParams",
      "definitions": {
        "EchoParams": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          },
          "additionalProperties": false,
          "required": [
            "name"
          ]
        }
      }
    },
    "result": {
      "$ref": "#/definitions/EchoResult",
      "definitions": {
        "EchoResult": {
          "type": "object",
          "properties": {
            "message": {
              "type": "string"
            }
          },
          "additionalProperties": false,
          "required": [
            "message"
          ]
        }
      }
    }
  }
]
```

## License

Released under the [MIT License](https://github.com/osamingo/jsonrpc/blob/master/LICENSE).
