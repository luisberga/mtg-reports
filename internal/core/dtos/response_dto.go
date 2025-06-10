package dtos

import "time"

type ResponseInsertCard struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
	Foil            bool   `json:"foil"`
}

type ResponseCard struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Set             string    `json:"set"`
	CollectorNumber string    `json:"collector_number"`
	Foil            bool      `json:"foil"`
	LastPrice       float64   `json:"last_price"`
	OldPrice        float64   `json:"old_price"`
	PriceChange     float64   `json:"price_change"`
	LastUpdate      time.Time `json:"last_update"`
}

type ResponseConciliateJob struct {
	Processed    int64 `json:"processed"`
	NotProcessed int64 `json:"not_processed"`
}

type ResponsePaginatedCards struct {
	Cards      []ResponseCard `json:"cards"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type ResponseCollectionStats struct {
	TotalCards int64   `json:"total_cards"`
	FoilCards  int64   `json:"foil_cards"`
	UniqueSets int64   `json:"unique_sets"`
	TotalValue float64 `json:"total_value"`
}
