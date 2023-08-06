package web

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

type web struct {
	http *http.Client
}

func New() *web {
	return &web{
		http: &http.Client{
			Timeout: time.Second * 30, // include in config
		},
	}
}

func (c *web) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

type HTTPResponse struct {
	Resp *http.Response
}

func (r *HTTPResponse) Body() io.ReadCloser {
	return r.Resp.Body
}

func (r *HTTPResponse) StatusCode() int {
	return r.Resp.StatusCode
}

func (c *web) Do(req Request) (Response, error) {
	request, ok := req.(*http.Request)
	if !ok {
		return nil, errors.New("invalid http request")
	}

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	return &HTTPResponse{Resp: resp}, nil
}
