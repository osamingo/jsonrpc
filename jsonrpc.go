package jsonrpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	// Version is JSON-RPC 2.0.
	Version = "2.0"

	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json"
)

var (
	batchRequestPrefixKey = byte('[')
	batchRequestSuffixKey = byte(']')
)

type (
	// A Request represents a JSON-RPC request received by the server.
	Request struct {
		ID      string           `json:"id"`
		Version string           `json:"jsonrpc"`
		Method  string           `json:"method"`
		Params  *json.RawMessage `json:"params"`
	}

	// A Response represents a JSON-RPC response returned by the server.
	Response struct {
		ID      string      `json:"id"`
		Version string      `json:"jsonrpc"`
		Result  interface{} `json:"result,omitempty"`
		Error   *Error      `json:"error,omitempty"`
	}
)

// ParseRequest parses a HTTP request to JSON-RPC request.
func ParseRequest(r *http.Request) ([]Request, *Error) {

	if r.Header.Get(contentTypeKey) != contentTypeValue {
		return nil, ErrInvalidRequest()
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrInvalidRequest()
	}
	r.Body.Close()

	if len(body) == 0 {
		return nil, ErrInvalidRequest()
	}

	if body[0] != batchRequestPrefixKey {
		var req Request
		if err = json.Unmarshal(body, &req); err != nil {
			return nil, ErrParse()
		}
		return []Request{req}, nil
	}

	if body[len(body)-1] != batchRequestSuffixKey {
		return nil, ErrParse()
	}

	rs := []Request{}
	if err = json.Unmarshal(body, &rs); err != nil {
		return nil, ErrParse()
	}

	return rs, nil
}

// NewResponse generates a JSON-RPC response.
func NewResponse(r Request) Response {
	res := Response{
		Version: r.Version,
	}
	if r.ID != "" {
		res.ID = r.ID
	}
	return res
}

// SendResponse writes JSON-RPC response.
func SendResponse(w http.ResponseWriter, resp []Response) error {

	var bin []byte
	var err error
	if len(resp) == 1 {
		bin, err = json.Marshal(&resp[0])
	} else if len(resp) > 1 {
		bin, err = json.Marshal(&resp)
	}

	if err != nil {
		return err
	}

	w.Header().Set(contentTypeKey, contentTypeValue)
	w.Write(bin)

	return nil
}
