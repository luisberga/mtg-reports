package validate

import (
	"errors"
	"mtg-report/internal/core/dtos"
	"strconv"
)

type validator struct{}

func New() *validator {
	return &validator{}
}

func (v *validator) Card(card dtos.RequestInsertCard) error {
	if card.Name == "" {
		return errors.New("name is required")
	}

	if card.CollectorNumber == "" {
		return errors.New("collector_number is required")
	}

	if card.SetName == "" {
		return errors.New("set_name is required")
	}

	if card.Foil == nil {
		return errors.New("foil is required")
	}

	if *card.Foil != true && *card.Foil != false {
		return errors.New("foil must be true or false")
	}

	return nil
}

func (v *validator) CardID(parts []string) (string, error) {
	var id string

	if len(parts) != 3 {
		return "", errors.New("invalid url")
	}

	if len(parts[2]) == 0 {
		return id, errors.New("id is required")
	}

	_, err := strconv.Atoi(parts[2])
	if err != nil {
		return id, errors.New("invalid id")
	}

	id = parts[2]

	return id, nil
}

func (v *validator) CardName(card dtos.RequestUpdateCard) error {
	if card.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (v *validator) Filters(setName, name, collector_number string) map[string]string {
	filters := make(map[string]string)

	if len(setName) != 0 {
		filters["set_name"] = setName
	}

	if len(name) != 0 {
		filters["name"] = name
	}

	if len(collector_number) != 0 {
		filters["collector_number"] = collector_number
	}

	return filters
}

func (v *validator) Pagination(pageStr, limitStr string) (int, int, error) {
	page := 1
	limit := 20 // default limit

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, errors.New("invalid page parameter")
		}
		if p < 1 {
			return 0, 0, errors.New("page must be greater than 0")
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, errors.New("invalid limit parameter")
		}
		if l < 1 || l > 100 {
			return 0, 0, errors.New("limit must be between 1 and 100")
		}
		limit = l
	}

	return page, limit, nil
}
