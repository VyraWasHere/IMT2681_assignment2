package assignment2

import (
	"net/http"
)

func rateHandler(w http.ResponseWriter, r *http.Request) {
	var payload Webhook
	err := json.Unmarshal(r.Body, &payload)
}