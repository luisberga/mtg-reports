package cardgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/sources/logger/logrus"
	"mtg-report/internal/sources/web"
	"net/http"
	"strconv"
)

type cardGateway struct {
	web web.HTTP
	log logrus.Logger
}

func New(web web.HTTP, log logrus.Logger) *cardGateway {
	return &cardGateway{
		web: web,
		log: log,
	}
}

func (cg *cardGateway) GetCardPrice(ctx context.Context, card domain.Cards) (float64, error) {
	url := fmt.Sprintf("https://api.scryfall.com/cards/%s/%s", card.SetName, card.CollectorNumber)
	req, err := cg.web.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("card gateway failed to get card: %w", err)
	}

	resp, err := cg.web.Do(req)
	if err != nil {
		return 0, fmt.Errorf("card gateway failed to get response: %w", err)
	}
	defer resp.Body().Close()

	if resp.StatusCode() == http.StatusNotFound {
		return 0, ErrCardNotFound{}
	}

	body, err := ioutil.ReadAll(resp.Body())
	if err != nil {
		return 0, fmt.Errorf("card gateway failed to read body: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("card gateway failed to get card: http status %d, response: %s", resp.StatusCode(), string(body))
	}

	var cardRequest entities.ScryfallCard
	err = json.Unmarshal(body, &cardRequest)
	if err != nil {
		return 0, fmt.Errorf("card gateway failed to unmarshal body: %w", err)
	}

	if card.Foil {
		if cardRequest.Prices.USDFoil != nil {
			usdFoil, err := strconv.ParseFloat(*cardRequest.Prices.USDFoil, 64)
			if err != nil {
				return 0, fmt.Errorf("card gateway failed to parse float for usd foil: %w", err)
			}
			return usdFoil, nil
		}
	} else {
		if cardRequest.Prices.USD != nil {
			usd, err := strconv.ParseFloat(*cardRequest.Prices.USD, 64)
			if err != nil {
				return 0, fmt.Errorf("card gateway failed to parse float for usd: %w", err)
			}
			return usd, nil
		}
	}

	return 0, ErrPriceIsZero{}
}
