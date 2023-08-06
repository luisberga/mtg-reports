package mocks

import (
	"context"
	"io"
	"mtg-report/internal/sources/web"

	"github.com/stretchr/testify/mock"
)

type httpMock struct {
	mock.Mock
}

func NewHTTPMock() *httpMock {
	return &httpMock{}
}

func (h *httpMock) NewRequestWithContext(ctx context.Context, method string, url string, body interface{}) (web.Request, error) {
	argsMock := h.Called(ctx, method, url, body)
	return argsMock.Get(0).(web.Response), argsMock.Error(1)
}

func (h *httpMock) Do(req context.Context) (web.Response, error) {
	argsMock := h.Called(req)
	return argsMock.Get(0).(web.Response), argsMock.Error(1)
}

type responseMock struct {
	mock.Mock
}

func NewResponseMock() *responseMock {
	return &responseMock{}
}

func (r *responseMock) Body() io.ReadCloser {
	argsMock := r.Called()
	return argsMock.Get(0).(io.ReadCloser)
}

func (r *responseMock) StatusCode() int {
	argsMock := r.Called()
	return argsMock.Get(0).(int)
}

type requestMock struct {
	mock.Mock
}
