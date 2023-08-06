package reporthandler

import (
	"context"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
)

type handler struct {
	ReportService ports.ReportService
	log           logrus.Logger
}

func New(cr ports.ReportService, log logrus.Logger) *handler {
	return &handler{
		ReportService: cr,
		log:           log,
	}
}

func (h *handler) ProcessAndSend(ctx context.Context) error {
	h.log.Info("process and send")

	err := h.ReportService.ProcessAndSend(ctx)
	if err != nil {
		return err
	}

	h.log.Info("email sent")

	return nil
}
