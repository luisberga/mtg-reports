package dtos

type RequestInsertCard struct {
	Name            string `json:"name,omitempty"`
	SetName         string `json:"set_name,omitempty"`
	CollectorNumber string `json:"collector_number,omitempty"`
	Foil            *bool  `json:"foil,omitempty"`
}

type RequestUpdateCard struct {
	ID   string
	Name string `json:"name,omitempty"`
}
