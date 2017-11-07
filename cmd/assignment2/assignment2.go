package main

import (
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/exchange/", exchangeHandler)
	http.ListenAndServe(":8080", nil)
}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	switch parts[2] {
	case "":
		if r.Method == http.MethodPost {
			webhookHandler(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

	case "avarage":
		avgHandler(w, r)

	case "latest":
		rateHandler(w, r)

	default:
		if r.Method == http.MethodGet {
			webhookHandler(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

}
