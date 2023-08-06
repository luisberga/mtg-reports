package conciliatehandler

import (
	"context"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
)

type handler struct {
	PriceService ports.PriceService
	log          logrus.Logger
}

func New(cs ports.PriceService, log logrus.Logger) *handler {
	return &handler{
		PriceService: cs,
		log:          log,
	}
}

func (h *handler) Conciliate(ctx context.Context) error {
	h.log.Info("conciliate")

	cardsUpdated, err := h.PriceService.Conciliate(ctx)
	if err != nil {
		h.log.WithError(err).Error("failed to conciliate")
	}

	h.log.WithFields(logrus.Fields{
		"cards_updated": cardsUpdated,
	}).Info("job done")

	return nil
}
