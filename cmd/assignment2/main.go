package main

import (
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/exchange", exchangeHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 2 {
		if r.Method == http.MethodPost {
			webhookHandler(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else if len(parts) > 2 {
		switch parts[2] {

		case "avarage":
			avgHandler(w, r)

		case "latest":
			rateHandler(w, r)

		case "evaluationtrigger":
			evalHandler(w, r)

		case "":
			if r.Method == http.MethodPost {
				webhookHandler(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

		default:
			if r.Method == http.MethodGet || r.Method == http.MethodDelete {
				webhookHandler(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
