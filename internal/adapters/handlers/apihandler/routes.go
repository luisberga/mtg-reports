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

	mux.HandleFunc("/card", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			c.InsertCard(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/card/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.GetCardbyID(w, r)
		case http.MethodPatch:
			c.UpdateCard(w, r)
		case http.MethodDelete:
			c.DeleteCard(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			c.InsertCards(w, r)
		case http.MethodGet:
			c.GetCards(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/card-history/", c.GetCardHistory)

	return mux
}
