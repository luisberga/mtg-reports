package apihandler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/dtos"
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

	assert.NotNil(t, h)
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
		{
			name:      "should return StatusInternalServerError when unable to read request body",
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
		{
			name:      "should return StatusBadRequest when unable to unmarshal request body",
			reqMethod: http.MethodPost,
			reqBody:   []byte("{invalid json}"),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "should return StatusBadRequest when validation fails",
			reqMethod: http.MethodPost,
			reqBody:   []byte(`{"name": ""}`),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("Card", mock.Anything).Return(errors.New("name is required"))
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "should return StatusBadRequest when card already exists",
			reqMethod: http.MethodPost,
			reqBody:   []byte(`{"name": "Card1", "set_name": "M21", "collector_number": "123", "foil": true}`),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("Card", mock.Anything).Return(nil)
				sMock.On("InsertCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{}, domain.ErrCardAlreadyExists{})
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "should return StatusOK when insert is successful",
			reqMethod: http.MethodPost,
			reqBody:   []byte(`{"name": "Card1", "set_name": "M21", "collector_number": "123", "foil": true}`),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("Card", mock.Anything).Return(nil)
				sMock.On("InsertCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{ID: 1, Name: "Card1"}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
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

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_InsertCards(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func() (*bytes.Buffer, string)
		mockSetup func(
			sMock *mocks.CardServiceMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusOK when cards are inserted successfully",
			setupFile: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("cards", "cards.csv")
				part.Write([]byte("name,set,number\nCard1,M21,123"))
				contentType := writer.FormDataContentType()
				writer.Close()
				return body, contentType
			},
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				sMock.On("InsertCards", mock.Anything, mock.Anything).Return(int64(1), int64(0))
			},
			wantCode: http.StatusOK,
		},
		{
			name: "should return StatusInternalServerError when no file is provided",
			setupFile: func() (*bytes.Buffer, string) {
				return &bytes.Buffer{}, "application/json"
			},
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			body, contentType := tt.setupFile()
			req, _ := http.NewRequest(http.MethodPost, "/cards", body)
			req.Header.Set("Content-Type", contentType)

			resp := httptest.NewRecorder()
			h.InsertCards(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_GetCardbyID(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusBadRequest when validation fails",
			url:  "/card/invalid",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("", errors.New("invalid id"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusBadRequest when card not found",
			url:  "/card/999",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("999", nil)
				sMock.On("GetCardbyID", mock.Anything, "999").Return(dtos.ResponseCard{}, domain.ErrCardNotFound{})
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusOK when card is found",
			url:  "/card/1",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("CardID", mock.Anything).Return("1", nil)
				sMock.On("GetCardbyID", mock.Anything, "1").Return(dtos.ResponseCard{ID: 1, Name: "Card1"}, nil)
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			resp := httptest.NewRecorder()

			h.GetCardbyID(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_GetCards(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusBadRequest when pagination validation fails",
			url:  "/cards?page=invalid&limit=10",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("Filters", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{})
				vMock.On("Pagination", "invalid", "10").Return(0, 0, errors.New("invalid page"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusOK when cards are retrieved",
			url:  "/cards?page=1&limit=10",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("Filters", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{})
				vMock.On("Pagination", "1", "10").Return(1, 10, nil)
				sMock.On("GetCardsPaginated", mock.Anything, mock.Anything, 1, 10).Return(dtos.ResponsePaginatedCards{}, nil)
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			resp := httptest.NewRecorder()

			h.GetCards(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_DeleteCard(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusBadRequest when validation fails",
			url:  "/card/invalid",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("", errors.New("invalid id"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusOK when card is deleted",
			url:  "/card/1",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("CardID", mock.Anything).Return("1", nil)
				sMock.On("DeleteCard", mock.Anything, "1").Return(nil)
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			req, _ := http.NewRequest(http.MethodDelete, tt.url, nil)
			resp := httptest.NewRecorder()

			h.DeleteCard(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_GetCardHistory(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusBadRequest when validation fails",
			url:  "/card/invalid/history",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("", errors.New("invalid id"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusBadRequest when pagination validation fails",
			url:  "/card/1/history?page=invalid&limit=10",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("1", nil)
				vMock.On("Pagination", "invalid", "10").Return(0, 0, errors.New("invalid page"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "should return StatusOK when history is retrieved",
			url:  "/card/1/history?page=1&limit=10",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("CardID", mock.Anything).Return("1", nil)
				vMock.On("Pagination", "1", "10").Return(1, 10, nil)
				sMock.On("GetCardHistoryPaginated", mock.Anything, "1", 1, 10).Return(dtos.ResponsePaginatedCards{}, nil)
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			resp := httptest.NewRecorder()

			h.GetCardHistory(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_GetCollectionStats(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name: "should return StatusOK when stats are retrieved",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				sMock.On("GetCollectionStats", mock.Anything).Return(dtos.ResponseCollectionStats{
					TotalCards: 100,
					FoilCards:  25,
					UniqueSets: 10,
					TotalValue: 1500.50,
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "should return StatusInternalServerError when service fails",
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Error", mock.Anything).Once()
				sMock.On("GetCollectionStats", mock.Anything).Return(dtos.ResponseCollectionStats{}, errors.New("service error"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sMock := mocks.NewCardServiceMock()
			vMock := mocks.NewValidateMock()
			lMock := mocks.NewLogMock()
			cMock := mocks.NewCustomMock()

			tt.mockSetup(sMock, vMock, lMock, cMock)

			h := New(vMock, sMock, lMock)

			req, _ := http.NewRequest(http.MethodGet, "/collection-stats", nil)
			resp := httptest.NewRecorder()

			h.GetCollectionStats(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_UpdateCard(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		reqBody   interface{}
		mockSetup func(
			sMock *mocks.CardServiceMock,
			vMock *mocks.ValidateMock,
			lMock *mocks.LogMock,
			cMock *mocks.CustomMock,
		)
		wantCode int
	}{
		{
			name:    "should return StatusBadRequest when validation fails",
			url:     "/card/invalid",
			reqBody: []byte(`{"name": "Updated Card"}`),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("", errors.New("invalid id"))
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "should return StatusInternalServerError when unable to read body",
			url:     "/card/1",
			reqBody: errorReader{},
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Once()
				lMock.On("WithError", mock.Anything).Return(cMock).Once()
				cMock.On("Warn", mock.Anything).Once()
				vMock.On("CardID", mock.Anything).Return("1", nil)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:    "should return StatusOK when card is updated",
			url:     "/card/1",
			reqBody: []byte(`{"name": "Updated Card"}`),
			mockSetup: func(
				sMock *mocks.CardServiceMock,
				vMock *mocks.ValidateMock,
				lMock *mocks.LogMock,
				cMock *mocks.CustomMock,
			) {
				lMock.On("Info", mock.Anything).Twice()
				vMock.On("CardID", mock.Anything).Return("1", nil)
				vMock.On("CardName", mock.Anything).Return(nil)
				sMock.On("UpdateCard", mock.Anything, mock.Anything).Return(dtos.ResponseInsertCard{ID: 1, Name: "Updated Card"}, nil)
			},
			wantCode: http.StatusOK,
		},
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

			req, _ := http.NewRequest(http.MethodPatch, tt.url, body)
			resp := httptest.NewRecorder()

			h.UpdateCard(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			sMock.AssertExpectations(t)
			vMock.AssertExpectations(t)
			lMock.AssertExpectations(t)
			cMock.AssertExpectations(t)
		})
	}
}

func Test_encondeResponse(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
		wantCode int
	}{
		{
			name:     "should encode response successfully",
			response: map[string]string{"test": "value"},
			wantCode: http.StatusOK,
		},
		{
			name:     "should handle nil response",
			response: nil,
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			encondeResponse(resp, tt.response)
			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		})
	}
}
