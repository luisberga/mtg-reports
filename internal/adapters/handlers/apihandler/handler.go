package apihandler

import (
	"encoding/json"
	"errors"
	"io"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/dtos"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
	"net/http"
	"strings"
)

type validate interface {
	Card(dtos.RequestInsertCard) error
	CardID(parts []string) (string, error)
	Filters(setName, name, collector_number string) map[string]string
	CardName(card dtos.RequestUpdateCard) error
	Pagination(pageStr, limitStr string) (int, int, error)
}

type apiHandler struct {
	validator   validate
	CardService ports.CardService
	log         logrus.Logger
}

func New(v validate, cs ports.CardService, log logrus.Logger) *apiHandler {
	return &apiHandler{
		validator:   v,
		CardService: cs,
		log:         log,
	}
}

func (h *apiHandler) InsertCard(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler insert card")

	card := dtos.RequestInsertCard{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.WithError(err).Warn("error to read body on insert card")
		http.Error(w, "failed to insert card", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &card)
	if err != nil {
		h.log.WithError(err).Warn("error to read body on insert card")
		http.Error(w, "failed to insert card, check body", http.StatusBadRequest)
		return
	}

	err = h.validator.Card(card)
	if err != nil {
		h.log.WithError(err).Warn("failed to insert card")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.CardService.InsertCard(r.Context(), card)
	if errors.Is(err, domain.ErrCardAlreadyExists{}) {
		h.log.WithError(err).Warn("failed to insert card")
		http.Error(w, domain.ErrCardAlreadyExists{}.Error(), http.StatusBadRequest)
	} else if errors.Is(err, domain.ErrInvalidSetName{}) {
		h.log.WithError(err).Warn("failed to insert card")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if err != nil {
		h.log.WithError(err).Error("failed to insert card")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("card inserted")
		encondeResponse(w, response)
	}
}

// TODO use queue here in future
func (h *apiHandler) InsertCards(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler insert cards")

	file, _, err := r.FormFile("cards")
	if err != nil {
		h.log.WithError(err).Warn("failed to insert card")
		http.Error(w, "failed to receive file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	inserted, failed := h.CardService.InsertCards(r.Context(), file)

	response := dtos.ResponseConciliateJob{
		Processed:    inserted,
		NotProcessed: failed,
	}

	h.log.Info("cards inserted")
	encondeResponse(w, response)
}

func (h *apiHandler) GetCardbyID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler get card by id")

	parts := strings.Split(r.URL.Path, "/")
	id, err := h.validator.CardID(parts)
	if err != nil {
		h.log.WithError(err).Warn("failed to get card by id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.CardService.GetCardbyID(r.Context(), id)
	if errors.Is(err, domain.ErrCardNotFound{}) {
		h.log.WithError(err).Warn("failed to get card by id")
		http.Error(w, domain.ErrCardNotFound{}.Error(), http.StatusBadRequest)
	} else if err != nil {
		h.log.WithError(err).Error("failed to get card by id")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("card retrieved")
		encondeResponse(w, response)
	}
}

func (h *apiHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler get cards")

	setName := r.URL.Query().Get("set_name")
	name := r.URL.Query().Get("name")
	collector_number := r.URL.Query().Get("collector_number")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	filters := h.validator.Filters(setName, name, collector_number)

	page, limit, err := h.validator.Pagination(pageStr, limitStr)
	if err != nil {
		h.log.WithError(err).Warn("failed to validate pagination parameters")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.CardService.GetCardsPaginated(r.Context(), filters, page, limit)
	if err != nil {
		h.log.WithError(err).Error("failed to get cards paginated")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("cards retrieved")
		encondeResponse(w, response)
	}
}

func (h *apiHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler delete cards")

	parts := strings.Split(r.URL.Path, "/")
	id, err := h.validator.CardID(parts)
	if err != nil {
		h.log.WithError(err).Warn("failed to get card by id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.CardService.DeleteCard(r.Context(), id)
	if err != nil {
		h.log.WithError(err).Error("failed to delete cards")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("card deleted")
	}
}

func (h *apiHandler) GetCardHistory(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler get card history")

	parts := strings.Split(r.URL.Path, "/")
	id, err := h.validator.CardID(parts)
	if err != nil {
		h.log.WithError(err).Warn("failed to get card history")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := h.validator.Pagination(pageStr, limitStr)
	if err != nil {
		h.log.WithError(err).Warn("failed to validate pagination parameters")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.CardService.GetCardHistoryPaginated(r.Context(), id, page, limit)
	if errors.Is(err, domain.ErrCardNotFound{}) {
		h.log.WithError(err).Warn("failed to get card history")
		http.Error(w, domain.ErrCardNotFound{}.Error(), http.StatusBadRequest)
	} else if err != nil {
		h.log.WithError(err).Error("failed to get card history")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("card history retrieved")
		encondeResponse(w, response)
	}
}

func (h *apiHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler update card")

	parts := strings.Split(r.URL.Path, "/")
	id, err := h.validator.CardID(parts)
	if err != nil {
		h.log.WithError(err).Warn("failed to update card")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card := dtos.RequestUpdateCard{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.WithError(err).Warn("error to read body on update card")
		http.Error(w, "failed to update card", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &card)
	if err != nil {
		h.log.WithError(err).Warn("error to read body on update card")
		http.Error(w, "failed to update card, check body", http.StatusBadRequest)
		return
	}

	err = h.validator.CardName(card)
	if err != nil {
		h.log.WithError(err).Warn("failed to update card")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card.ID = id

	response, err := h.CardService.UpdateCard(r.Context(), card)
	if errors.Is(err, domain.ErrCardNotFound{}) {
		h.log.WithError(err).Warn("failed to update card")
		http.Error(w, domain.ErrCardNotFound{}.Error(), http.StatusBadRequest)
	} else if err != nil {
		h.log.WithError(err).Error("failed to update card")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("card updated")
		encondeResponse(w, response)
	}
}

func (h *apiHandler) GetCollectionStats(w http.ResponseWriter, r *http.Request) {
	h.log.Info("handler get collection stats")

	response, err := h.CardService.GetCollectionStats(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get collection stats")
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
	} else {
		h.log.Info("collection stats retrieved")
		encondeResponse(w, response)
	}
}

func encondeResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, ErrInternalErr{}.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
