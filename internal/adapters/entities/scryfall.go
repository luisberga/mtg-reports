package entities

type Price struct {
	USD     *string `json:"usd"`
	USDFoil *string `json:"usd_foil"`
}

type ScryfallCard struct {
	Prices Price `json:"prices"`
}
