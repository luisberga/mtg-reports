package exchangegateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/sources/logger/logrus"
	"mtg-report/internal/sources/web"
	"net/http"
)

type exchangeGateway struct {
	web web.HTTP
	log logrus.Logger
	url string
}

func New(web web.HTTP, url string, log logrus.Logger) *exchangeGateway {
	return &exchangeGateway{
		web: web,
		log: log,
		url: url,
	}
}

func (eg *exchangeGateway) GetUSD(ctx context.Context) (float64, error) {
	url := fmt.Sprintf(eg.url)
	req, err := eg.web.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("exchange gateway failed to create request: %w", err)
	}

	resp, err := eg.web.Do(req)
	if err != nil {
		return 0, fmt.Errorf("exchange gateway failed to get response: %w", err)
	}
	defer resp.Body().Close()

	if resp.StatusCode() != http.StatusOK {
		return 0, ErrFailedToGetExchangeRequest{}
	}

	body, err := ioutil.ReadAll(resp.Body())
	if err != nil {
		return 0, fmt.Errorf("exchange gateway failed to read body: %w", err)
	}

	var exchange entities.ExchangeRate
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return 0, fmt.Errorf("exchange gateway failed to unmarshal body: %w", err)
	}

	if exchange.ConversionRates.BRL == nil {
		return 0, ErrExchangeRequestNillValue{}
	}

	return *exchange.ConversionRates.BRL, nil
}
