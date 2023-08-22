package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"istomyang.github.com/like-iam/component-base/base"
	"net/http"
)

type response struct {
	native  *http.Response
	decoder base.Decoder

	body []byte
	err  error
}

func newResponse(res *http.Response, decoder base.Decoder, err error) Response {
	var r response
	r.native = res
	r.decoder = decoder
	r.err = err
	return &r
}

func (r *response) Into(v any) error {
	body, err := r.getBody()
	if err != nil {
		return err
	}
	return r.decoder.Decode(body, v)
}

func (r *response) Raw() ([]byte, error) {
	return r.getBody()
}

func (r *response) Error() error {
	// r.err will not be nil when have error in build request stage, which r.native is nil,
	// so err doesn't need to be slice type.
	if r.err != nil {
		return r.err
	}
	if r.native.StatusCode > 500 {
		if body, err := r.getBody(); err != nil {
			return err
		} else {
			return errors.New(string(body))
		}
	}
	return nil
}

func (r *response) getBody() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	var out bytes.Buffer
	var body = r.native.Body
	if body == nil {
		return nil, fmt.Errorf("response body is nil: %v", r.native)
	}
	if _, err := io.Copy(&out, body); err != nil {
		return nil, err
	}
	if err := body.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

var _ Response = &response{}
