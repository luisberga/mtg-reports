package reportservice

import (
	"context"
	"fmt"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
	"strings"
	"time"
)

type service struct {
	ReportRepository ports.ReportRepository
	Email            ports.Email
	log              logrus.Logger
}

func New(rr ports.ReportRepository, email ports.Email, log logrus.Logger) *service {
	return &service{
		ReportRepository: rr,
		Email:            email,
		log:              log,
	}
}

func (s *service) ProcessAndSend(ctx context.Context) error {
	err := s.ReportRepository.InsertTotalPrice(ctx)
	if err != nil {
		return fmt.Errorf("service failed to insert total price in process and send: %w", err)
	}

	cards, err := s.ReportRepository.GetCardsReport(ctx)
	if err != nil {
		return fmt.Errorf("service failed to get cards reports in process and send: %w", err)
	}

	cardsTable := s.formatCardsTable(cards)

	cardsPrice, err := s.ReportRepository.GetTotalPrice(ctx)
	if err != nil {
		return fmt.Errorf("service failed to get total price in process and send: %w", err)
	}

	cardsPriceFormatted := s.formatCardsPrice(cardsPrice)

	err = s.Email.SendEmail(cardsTable, cardsPriceFormatted)
	if err != nil {
		return fmt.Errorf("service failed to send email in process and send: %w", err)
	}

	return nil
}

func (s *service) formatCardsTable(cards []domain.Cards) string {
	var builder strings.Builder

	header := "<tr>" +
		"<th style='border: 1px solid black; padding: 10px;'>ID</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Name</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Set Name</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Collector Number</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Foil</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Old Price</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Last Price</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Price Change</th>" +
		"<th style='border: 1px solid black; padding: 10px;'>Last Update</th>" +
		"</tr>"
	builder.WriteString(fmt.Sprintf(header))

	rowFormat := "<tr>" +
		"<td style='border: 1px solid black; padding: 10px;'>%v</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%s</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%s</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%s</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%v</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%.2f</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%.2f</td>" +
		"<td style='border: 1px solid black; padding: 10px; color: %s;'>%.2f</td>" +
		"<td style='border: 1px solid black; padding: 10px;'>%s</td>" +
		"</tr>"

	for _, card := range cards {
		color := "black"
		if card.PriceChange > 0 {
			color = "green"
		} else if card.PriceChange < 0 {
			color = "red"
		}

		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		row := fmt.Sprintf(rowFormat,
			card.ID, card.Name, card.SetName, card.CollectorNumber,
			card.Foil, card.OldPrice, card.LastPrice, color, card.PriceChange,
			lastUpdate.Format(time.RFC1123))
		builder.WriteString(row)
	}

	builder.WriteString("</table>")

	return builder.String()
}

func (s *service) formatCardsPrice(price domain.CardsPrice) string {
	var builder strings.Builder

	var color string
	if price.PriceChange > 0 {
		color = "<span style='color: green;'>increased</span>"
	} else if price.PriceChange < 0 {
		color = "<span style='color: red;'>decreased</span>"
	} else {
		color = "<span style='color: black;'>stayed the same</span>"
	}

	builder.WriteString(fmt.Sprintf("The total value of your MTG card investment has %s from <strong>$%.2f</strong> to <strong>$%.2f</strong>. That's a change of <strong>$%.2f</strong>.",
		color, price.OldPrice, price.NewPrice, price.PriceChange))

	return builder.String()
}
