package entities

import "time"

type MysqlCardInfo struct {
	ID              int64    `db:"id"`
	Name            string   `db:"name"`
	SetName         string   `db:"set_name"`
	CollectorNumber string   `db:"collector_number"`
	LastPrice       *float64 `db:"last_price"`
	Foil            bool     `db:"foil"`
}

type MysqlCardPriceHistory struct {
	ID              int64      `db:"id"`
	Name            string     `db:"name"`
	SetName         string     `db:"set_name"`
	CollectorNumber string     `db:"collector_number"`
	OldPrice        float64    `db:"last_price"`
	LastPrice       float64    `db:"last_price"`
	PriceChange     float64    `db:"price_change"`
	LastUpdate      *time.Time `db:"last_update"`
	Foil            bool       `db:"foil"`
}
