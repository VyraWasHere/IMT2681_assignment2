package assignment2

import (
	"net/http"
)

func init() {
	http.HandleFunc("/exchange", webhookHandler)
	http.HandleFunc("/exchange/avarage", avgHandler)
	http.HandleFunc("/exchange/latest", rateHandler)
	http.ListenAndServe(":8080", nil)
}
