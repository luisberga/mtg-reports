package apihandler

import "net/http"

type cards interface {
	InsertCard(http.ResponseWriter, *http.Request)
	InsertCards(http.ResponseWriter, *http.Request)
	GetCardbyID(http.ResponseWriter, *http.Request)
	GetCards(http.ResponseWriter, *http.Request)
	DeleteCard(http.ResponseWriter, *http.Request)
	GetCardHistory(w http.ResponseWriter, r *http.Request)
	UpdateCard(w http.ResponseWriter, r *http.Request)
}

func SetupRouter(c cards) *http.ServeMux {
	mux := http.NewServeMux()

	//POST
	mux.HandleFunc("/insert-card", c.InsertCard)

	//POST
	mux.HandleFunc("/insert-cards", c.InsertCards)

	//GET card/{id}
	mux.HandleFunc("/card/", c.GetCardbyID)

	//GET cards?set_name={set_name}?name={name}?collector_number={collector_number}
	mux.HandleFunc("/cards", c.GetCards)

	//DELETE delete-card/{id}
	mux.HandleFunc("/delete-card/", c.DeleteCard)

	//GET card-history/{id}
	mux.HandleFunc("/card-history/", c.GetCardHistory)

	//PATCH update-card/{id}
	mux.HandleFunc("/update-card/", c.UpdateCard)

	return mux
}
