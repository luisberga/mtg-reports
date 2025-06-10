package cardservice

import (
	"bufio"
	"context"
	"fmt"
	"mime/multipart"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/dtos"
	"mtg-report/internal/core/ports"
	"mtg-report/internal/sources/logger/logrus"
	"regexp"
	"strconv"
	"time"
)

type service struct {
	cardsRepository ports.CardsRepository
	commitSize      int
	log             logrus.Logger
}

func New(cr ports.CardsRepository, commitSize int, log logrus.Logger) *service {
	return &service{
		cardsRepository: cr,
		commitSize:      commitSize,
		log:             log,
	}
}

func (c *service) InsertCard(ctx context.Context, cardRequest dtos.RequestInsertCard) (dtos.ResponseInsertCard, error) {
	cardDomain := domain.Cards{
		Name:            cardRequest.Name,
		SetName:         cardRequest.SetName,
		CollectorNumber: cardRequest.CollectorNumber,
		Foil:            *cardRequest.Foil,
	}

	cardDomain, err := c.cardsRepository.InsertCard(ctx, cardDomain)
	if err != nil {
		return dtos.ResponseInsertCard{}, fmt.Errorf("service failed to insert card: %w", err)
	}

	return dtos.ResponseInsertCard{
		ID:              cardDomain.ID,
		Name:            cardDomain.Name,
		Set:             cardDomain.SetName,
		CollectorNumber: cardDomain.CollectorNumber,
		Foil:            cardDomain.Foil,
	}, nil
}

func (c *service) GetCardbyID(ctx context.Context, id string) (dtos.ResponseCard, error) {
	cardDomain, err := c.cardsRepository.GetCardbyID(ctx, id)
	if err != nil {
		return dtos.ResponseCard{}, fmt.Errorf("service failed to get card: %w", err)
	}

	var lastUpdate time.Time
	if cardDomain.LastUpdate != nil {
		lastUpdate = *cardDomain.LastUpdate
	}

	return dtos.ResponseCard{
		ID:              cardDomain.ID,
		Name:            cardDomain.Name,
		Set:             cardDomain.SetName,
		CollectorNumber: cardDomain.CollectorNumber,
		Foil:            cardDomain.Foil,
		LastPrice:       cardDomain.LastPrice,
		OldPrice:        cardDomain.OldPrice,
		PriceChange:     cardDomain.PriceChange,
		LastUpdate:      lastUpdate,
	}, nil
}

func (c *service) GetCards(ctx context.Context, filters map[string]string) ([]dtos.ResponseCard, error) {
	cardsDomain, err := c.cardsRepository.GetCards(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("service failed to get card: %w", err)
	}

	cards := make([]dtos.ResponseCard, 0, len(cardsDomain))
	for _, card := range cardsDomain {
		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		cards = append(cards, dtos.ResponseCard{
			ID:              card.ID,
			Name:            card.Name,
			Set:             card.SetName,
			CollectorNumber: card.CollectorNumber,
			Foil:            card.Foil,
			LastPrice:       card.LastPrice,
			OldPrice:        card.OldPrice,
			PriceChange:     card.PriceChange,
			LastUpdate:      lastUpdate,
		})
	}

	return cards, nil
}

func (c *service) UpdateCard(ctx context.Context, cardRequest dtos.RequestUpdateCard) (dtos.ResponseInsertCard, error) {
	id, err := strconv.ParseInt(cardRequest.ID, 10, 64)
	if err != nil {
		return dtos.ResponseInsertCard{}, fmt.Errorf("service failed to parse id in update card: %w", err)
	}

	updateCard := domain.UpdateCard{
		ID:   id,
		Name: cardRequest.Name,
	}

	cardsDomain, err := c.cardsRepository.UpdateCard(ctx, updateCard)
	if err != nil {
		return dtos.ResponseInsertCard{}, fmt.Errorf("service failed to update card: %w", err)
	}

	card := dtos.ResponseInsertCard{
		ID:              cardsDomain.ID,
		Name:            cardsDomain.Name,
		Set:             cardsDomain.SetName,
		CollectorNumber: cardsDomain.CollectorNumber,
		Foil:            cardsDomain.Foil,
	}

	return card, nil
}

func (c *service) DeleteCard(ctx context.Context, id string) error {
	err := c.cardsRepository.DeleteCard(ctx, id)
	if err != nil {
		return fmt.Errorf("service failed to delete card: %w", err)
	}

	return nil
}

func (c *service) GetCardHistory(ctx context.Context, id string) ([]dtos.ResponseCard, error) {
	cards, err := c.cardsRepository.GetCardHistory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service failed to get card history: %w", err)
	}

	cardsResponse := make([]dtos.ResponseCard, 0, len(cards))
	for _, card := range cards {
		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		cardsResponse = append(cardsResponse, dtos.ResponseCard{
			ID:              card.ID,
			Name:            card.Name,
			Set:             card.SetName,
			CollectorNumber: card.CollectorNumber,
			Foil:            card.Foil,
			LastPrice:       card.LastPrice,
			OldPrice:        card.OldPrice,
			PriceChange:     card.PriceChange,
			LastUpdate:      lastUpdate,
		})
	}

	return cardsResponse, nil
}

func (c *service) InsertCards(ctx context.Context, file multipart.File) (int64, int64) {
	var cardsProcessed int64
	var cardsNotProcessed int64
	cardsCh := make(chan []domain.Cards, 0)
	finishCh := make(chan struct{})

	scanner := bufio.NewScanner(file)

	re := regexp.MustCompile(`name: ([\p{L}\s-,'"!?]+), set_name: ([\p{L}\s-]+), collector_number: ([\w\s]+), foil: ([\w\s]+)`)

	go func() {
		defer close(cardsCh)
		cards := make([]domain.Cards, 0)

		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				continue
			}

			matches := re.FindStringSubmatch(line)
			if len(matches) != 5 {
				c.log.WithFields(logrus.Fields{"line": line}).Warn("service failed to parse line in insert cards")
				cardsNotProcessed++
				continue
			}

			card := domain.Cards{
				Name:            matches[1],
				SetName:         matches[2],
				CollectorNumber: matches[3],
			}

			if err := card.ValidateCardFields(matches[4]); err != nil {
				cardsNotProcessed++
				c.log.Warn(fmt.Errorf("service failed to insert one card in insert cards: %w", err))
				continue
			}

			cards = append(cards, card)
			if len(cards) == c.commitSize {
				cardsCh <- cards
				cards = make([]domain.Cards, 0, c.commitSize)
			}
		}

		if len(cards) > 0 {
			cardsCh <- cards
		}

	}()

	go func() {
		defer close(finishCh)

		for cards, ok := <-cardsCh; ok; cards, ok = <-cardsCh {
			err := c.cardsRepository.InsertCards(ctx, cards)
			if err != nil {
				c.log.Warn(fmt.Errorf("service failed to insert cards: %w", err))
				cardsNotProcessed += int64(len(cards))
				continue
			}
			cardsProcessed += int64(len(cards))
		}
	}()

	<-finishCh

	if err := scanner.Err(); err != nil {
		c.log.Error(fmt.Errorf("service scanner failed to insert cards: %w", err))
	}

	return cardsProcessed, cardsNotProcessed
}

func (c *service) GetCardsPaginated(ctx context.Context, filters map[string]string, page, limit int) (dtos.ResponsePaginatedCards, error) {
	offset := (page - 1) * limit

	// Get total count
	total, err := c.cardsRepository.GetCardsCount(ctx, filters)
	if err != nil {
		return dtos.ResponsePaginatedCards{}, fmt.Errorf("service failed to get cards count: %w", err)
	}

	// Get paginated cards
	cardsDomain, err := c.cardsRepository.GetCardsPaginated(ctx, filters, offset, limit)
	if err != nil {
		return dtos.ResponsePaginatedCards{}, fmt.Errorf("service failed to get cards paginated: %w", err)
	}

	cards := make([]dtos.ResponseCard, 0, len(cardsDomain))
	for _, card := range cardsDomain {
		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		cards = append(cards, dtos.ResponseCard{
			ID:              card.ID,
			Name:            card.Name,
			Set:             card.SetName,
			CollectorNumber: card.CollectorNumber,
			Foil:            card.Foil,
			LastPrice:       card.LastPrice,
			OldPrice:        card.OldPrice,
			PriceChange:     card.PriceChange,
			LastUpdate:      lastUpdate,
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division

	return dtos.ResponsePaginatedCards{
		Cards:      cards,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (c *service) GetCardHistoryPaginated(ctx context.Context, id string, page, limit int) (dtos.ResponsePaginatedCards, error) {
	offset := (page - 1) * limit

	// Get total count
	total, err := c.cardsRepository.GetCardHistoryCount(ctx, id)
	if err != nil {
		return dtos.ResponsePaginatedCards{}, fmt.Errorf("service failed to get card history count: %w", err)
	}

	// Get paginated card history
	cardsDomain, err := c.cardsRepository.GetCardHistoryPaginated(ctx, id, offset, limit)
	if err != nil {
		return dtos.ResponsePaginatedCards{}, fmt.Errorf("service failed to get card history paginated: %w", err)
	}

	cards := make([]dtos.ResponseCard, 0, len(cardsDomain))
	for _, card := range cardsDomain {
		var lastUpdate time.Time
		if card.LastUpdate != nil {
			lastUpdate = *card.LastUpdate
		}

		cards = append(cards, dtos.ResponseCard{
			ID:              card.ID,
			Name:            card.Name,
			Set:             card.SetName,
			CollectorNumber: card.CollectorNumber,
			Foil:            card.Foil,
			LastPrice:       card.LastPrice,
			OldPrice:        card.OldPrice,
			PriceChange:     card.PriceChange,
			LastUpdate:      lastUpdate,
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division

	return dtos.ResponsePaginatedCards{
		Cards:      cards,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (c *service) GetCollectionStats(ctx context.Context) (dtos.ResponseCollectionStats, error) {
	stats, err := c.cardsRepository.GetCollectionStats(ctx)
	if err != nil {
		return dtos.ResponseCollectionStats{}, fmt.Errorf("service failed to get collection stats: %w", err)
	}

	return dtos.ResponseCollectionStats{
		TotalCards: stats.TotalCards,
		FoilCards:  stats.FoilCards,
		UniqueSets: stats.UniqueSets,
		TotalValue: stats.TotalValue,
	}, nil
}
