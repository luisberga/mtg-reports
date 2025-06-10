package domain

import (
	"errors"
	"time"
)

type Cards struct {
	ID              int64
	Name            string
	SetName         string
	CollectorNumber string
	Foil            bool
	CardsDetails
}

func (c *Cards) ValidateCardFields(foil string) error {
	if len(c.Name) == 0 {
		return errors.New("name is required")
	}

	if len(c.CollectorNumber) == 0 {
		return errors.New("collector number is required")
	}

	if foil == "true" {
		c.Foil = true
	} else if foil == "false" {
		c.Foil = false
	} else {
		return errors.New("foil bool is required")
	}

	return nil
}

type CardsDetails struct {
	CardID      int64
	LastPrice   float64
	OldPrice    float64
	PriceChange float64
	LastUpdate  *time.Time
}

type UpdateCard struct {
	ID   int64
	Name string
}

type CardsPrice struct {
	OldPrice    float64
	NewPrice    float64
	PriceChange float64
	LastUpdate  *time.Time
}

type CollectionStats struct {
	TotalCards int64
	FoilCards  int64
	UniqueSets int64
	TotalValue float64
}
