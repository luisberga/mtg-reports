package conciliateservice

import (
	"context"
	"errors"
	"fmt"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
	"time"
)

const (
	exchangeDefault      float64 = 4.80
	maxRequestsPerSecond         = 10
)

type service struct {
	ConciliateRepository ports.ConciliateRepository
	cardGateway          ports.CardGateway
	exchangegateway      ports.ExchangeGateway
	commitSize           int
	log                  logrus.Logger
}

func New(cr ports.ConciliateRepository, cg ports.CardGateway, eg ports.ExchangeGateway, commitSize int, log logrus.Logger) *service {
	return &service{
		ConciliateRepository: cr,
		cardGateway:          cg,
		exchangegateway:      eg,
		commitSize:           commitSize,
		log:                  log,
	}
}

func (c *service) Conciliate(ctx context.Context) (int64, error) {
	var cardsUpdated int64

	exchangeValue, err := c.exchangegateway.GetUSD(ctx)
	if err != nil {
		c.log.Error(fmt.Errorf("service failed to get usd exchange: %w", err))
		exchangeValue = exchangeDefault
	}

	cardCh := make(chan []domain.CardsDetails, 0)
	finishCh := make(chan struct{})

	offset := 0

	go func() {
		defer close(cardCh)
		for {
			cards, err := c.ConciliateRepository.GetCardsForUpdate(ctx, offset, c.commitSize)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					c.log.Error(fmt.Errorf("service failed to get cards for update due context timeout: %w", err))
					break
				}
				c.log.Error(fmt.Errorf("service failed to get cards for update: %w", err))
			}

			if len(cards) == 0 {
				break
			}

			cardsDetails := make([]domain.CardsDetails, 0, len(cards))
			ticker := time.NewTicker(time.Second / maxRequestsPerSecond)
			defer ticker.Stop()

			for i, card := range cards {
				price, err := c.cardGateway.GetCardPrice(ctx, card)
				<-ticker.C
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						c.logError(card, fmt.Errorf("service failed to get card price due context timeout: %w", err))
						break
					}
					c.logError(card, fmt.Errorf("service failed to get card price: %w", err))
					continue
				}

				cards[i].CardsDetails.CardID = card.ID
				cards[i].OldPrice = card.LastPrice
				cards[i].LastPrice = price * exchangeValue
				cards[i].PriceChange = cards[i].LastPrice - cards[i].OldPrice

				lastUpdate := time.Now()
				cards[i].CardsDetails.LastUpdate = &lastUpdate
			}

			for _, card := range cards {
				if card.LastUpdate != nil {
					cardsDetails = append(cardsDetails, card.CardsDetails)
				}
			}

			cardCh <- cardsDetails
			offset = offset + c.commitSize
		}
	}()

	go func() {
		defer close(finishCh)
		for cards, ok := <-cardCh; ok; cards, ok = <-cardCh {
			c.log.Info("inserting cards...")
			if len(cards) == 0 {
				c.log.Info("no cards to insert")
				continue
			}

			err = c.ConciliateRepository.InsertCardDetails(ctx, cards)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					c.log.Error(fmt.Errorf("service failed to insert card details: %w", err))
					break
				}
				c.log.Warn(fmt.Errorf("service failed to insert card details: %w", err))
				continue
			}
			cardsUpdated = cardsUpdated + int64(len(cards))
			c.log.Info("cards inserted!")
		}
	}()

	<-finishCh

	return cardsUpdated, nil
}

func (c *service) logError(card domain.Cards, err error) {
	c.log.WithFields(logrus.Fields{
		"card_id":          card.ID,
		"card_name":        card.Name,
		"set_name":         card.SetName,
		"collector_number": card.CollectorNumber,
		"foil":             card.Foil,
	}).Warn(err)
}
