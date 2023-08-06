package web

import (
	"context"
	"io"
)

type HTTP interface {
	NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (Request, error)
	Do(req Request) (Response, error)
}

type Request interface{}

type Response interface {
	Body() io.ReadCloser
	StatusCode() int
}
