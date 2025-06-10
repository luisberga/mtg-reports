package apihandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCardsHandler struct {
	mock.Mock
}

func (m *mockCardsHandler) InsertCard(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) InsertCards(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) GetCardbyID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) GetCardHistory(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func (m *mockCardsHandler) GetCollectionStats(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func TestSetupRouter_CardPOST(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodPost, "/card", strings.NewReader("{}"))
	resp := httptest.NewRecorder()

	mockHandler.On("InsertCard", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardMethodNotAllowed(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/card", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Contains(t, resp.Body.String(), "Method not allowed")
}

func TestSetupRouter_CardWithIDGET(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/card/123", nil)
	resp := httptest.NewRecorder()

	mockHandler.On("GetCardbyID", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardWithIDPATCH(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodPatch, "/card/123", strings.NewReader("{}"))
	resp := httptest.NewRecorder()

	mockHandler.On("UpdateCard", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardWithIDDELETE(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodDelete, "/card/123", nil)
	resp := httptest.NewRecorder()

	mockHandler.On("DeleteCard", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardWithIDMethodNotAllowed(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodHead, "/card/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Contains(t, resp.Body.String(), "Method not allowed")
}

func TestSetupRouter_CardsPOST(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader("{}"))
	resp := httptest.NewRecorder()

	mockHandler.On("InsertCards", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardsGET(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/cards", nil)
	resp := httptest.NewRecorder()

	mockHandler.On("GetCards", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CardsMethodNotAllowed(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodDelete, "/cards", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Contains(t, resp.Body.String(), "Method not allowed")
}

func TestSetupRouter_CardHistory(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/card-history/123", nil)
	resp := httptest.NewRecorder()

	mockHandler.On("GetCardHistory", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CollectionStatsGET(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/collection-stats", nil)
	resp := httptest.NewRecorder()

	mockHandler.On("GetCollectionStats", resp, req)

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockHandler.AssertExpectations(t)
}

func TestSetupRouter_CollectionStatsMethodNotAllowed(t *testing.T) {
	mockHandler := &mockCardsHandler{}
	router := SetupRouter(mockHandler)

	req := httptest.NewRequest(http.MethodPost, "/collection-stats", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Contains(t, resp.Body.String(), "Method not allowed")
}
