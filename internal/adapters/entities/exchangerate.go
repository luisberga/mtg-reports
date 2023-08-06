package entities

type ConversionRates struct {
	BRL *float64 `json:"BRL"`
}

type ExchangeRate struct {
	ConversionRates ConversionRates `json:"conversion_rates"`
}
