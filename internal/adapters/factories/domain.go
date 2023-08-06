package factories

import (
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/core/domain"
	"time"
)

func CardsInfoToCardsDomain(cards []entities.MysqlCardInfo) []domain.Cards {
	domainCards := make([]domain.Cards, 0, len(cards))

	for _, card := range cards {
		var lastPrice float64
		if card.LastPrice != nil {
			lastPrice = *card.LastPrice
		}

		domainCards = append(domainCards, domain.Cards{
			ID:              card.ID,
			Name:            card.Name,
			SetName:         card.SetName,
			CollectorNumber: card.CollectorNumber,
			CardsDetails: domain.CardsDetails{
				LastPrice: lastPrice,
			},
			Foil: card.Foil,
		})
	}

	return domainCards
}

func CardPriceHistoryToCardsDomain(cards []entities.MysqlCardPriceHistory) []domain.Cards {
	domainCards := make([]domain.Cards, 0, len(cards))

	for _, card := range cards {
		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		domainCards = append(domainCards, domain.Cards{
			ID:              card.ID,
			Name:            card.Name,
			SetName:         card.SetName,
			CollectorNumber: card.CollectorNumber,
			CardsDetails: domain.CardsDetails{
				LastPrice:   card.LastPrice,
				OldPrice:    card.OldPrice,
				PriceChange: card.PriceChange,
				LastUpdate:  &lastUpdate,
			},
			Foil: card.Foil,
		})
	}

	return domainCards
}
