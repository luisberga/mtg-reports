package apihandler

import (
	"bytes"
	"errors"
	"io"
	"mtg-report/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorReader struct{}

func (er errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

func Test_New(t *testing.T) {
	sMock := mocks.NewCardServiceMock()
	vMock := mocks.NewValidateMock()
	lMock := mocks.NewLogMock()

	h := New(vMock, sMock, lMock)

	if h == nil {
		t.Error("expected not nil")
	}
}

func Test_InsertCard(t *testing.T) {
	tests := []struct {
		name      string
		reqMethod string
		reqBody   interface{}
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantErr  bool
		wantCode int
	}{
		// {
		// 	name:      "it should return StatusMethodNotAllowed when method is not POST",
		// 	reqMethod: http.MethodGet,
		// 	reqBody:   []byte("{}"),
		// 	mockSetup: func(sMock *mocks.CardServiceMock, vMock *mocks.ValidateMock, lMock *mocks.LogMock) {
		// 		lMock.On("Info", mock.Anything).Once()
		// 	},
		// 	wantErr:  true,
		// 	wantCode: http.StatusMethodNotAllowed,
		// },
		{
			name:      "it should return StatusInternalServerError when unable to read request body",
			reqMethod: http.MethodPost,
			reqBody:   errorReader{},
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				cMock.On("Warn", mock.Anything).Once()

				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		// {
		// 	name:      "it should return StatusBadRequest when unable to unmarshal request body",
		// 	reqMethod: http.MethodPost,
		// 	reqBody:   []byte("{aa}"),
		// 	mockSetup: func(sMock *mocks.CardServiceMock, vMock *mocks.ValidateMock, lMock *mocks.LogMock) {
		// 		lMock.On("Info", mock.Anything).Once()
		// 		lMock.On("Warn", mock.Anything).Once()
		// 	},
		// 	wantErr:  true,
		// 	wantCode: http.StatusBadRequest,
		// },
		// {
		// 	name:      "it should return StatusBadRequest when card already exists",
		// 	reqMethod: http.MethodPost,
		// 	reqBody:   []byte(`{"name": "Card1"}`),
		// 	mockSetup: func(sMock *mocks.CardServiceMock, vMock *mocks.ValidateMock, lMock *mocks.LogMock) {
		// 		lMock.On("Info", mock.Anything).Twice()
		// 		lMock.On("Warn", mock.Anything).Once()
		// 		vMock.On("Card", mock.Anything).Return(nil)
		// 		sMock.On("InsertCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{}, domain.ErrCardAlreadyExists{})
		// 	},
		// 	wantErr:  true,
		// 	wantCode: http.StatusBadRequest,
		// },
		// {
		// 	name:      "it should return StatusInternalServerError when there's an unknown error",
		// 	reqMethod: http.MethodPost,
		// 	reqBody:   []byte(`{"name": "Card1"}`),
		// 	mockSetup: func(sMock *mocks.CardServiceMock, vMock *mocks.ValidateMock, lMock *mocks.LogMock) {
		// 		lMock.On("Info", mock.Anything).Twice()
		// 		lMock.On("Error", mock.Anything).Once()
		// 		vMock.On("Card", mock.Anything).Return(nil)
		// 		sMock.On("InsertCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{}, errors.New("unknown error"))
		// 	},
		// 	wantErr:  true,
		// 	wantCode: http.StatusInternalServerError,
		// },
		// {
		// 	name:      "it should return StatusOK when insert is successful",
		// 	reqMethod: http.MethodPost,
		// 	reqBody:   []byte(`{"name": "Card1"}`),
		// 	mockSetup: func(sMock *mocks.CardServiceMock, vMock *mocks.ValidateMock, lMock *mocks.LogMock) {
		// 		lMock.On("Info", mock.Anything).Twice()
		// 		vMock.On("Card", mock.Anything).Return(nil)
		// 		sMock.On("InsertCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{}, nil)
		// 	},
		// 	wantErr:  false,
		// 	wantCode: http.StatusOK,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			var body io.Reader
			switch v := tt.reqBody.(type) {
			case []byte:
				body = bytes.NewBuffer(v)
			case errorReader:
				body = v
			default:
				t.Fatalf("unsupported type for reqBody: %T", tt.reqBody)
			}

			req, _ := http.NewRequest(tt.reqMethod, "/insert", body)
			resp := httptest.NewRecorder()

			h.InsertCard(resp, req)

			if tt.wantErr {
				assert.Equal(t, tt.wantCode, resp.Code)
			} else {
				assert.Equal(t, tt.wantCode, resp.Code)
			}

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}
