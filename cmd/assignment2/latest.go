package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func rateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Please provide base and target currency", http.StatusBadRequest)
		return
	}

	var payload Webhook
	pl, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(pl, &payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rate, body, status := getRate("latest", payload.BaseCurrency, payload.TargetCurrency)
	if status == http.StatusOK {
		fmt.Fprint(w, rate)
	} else {
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}
}

func getRate(when, base, target string) (rate float32, body string, status int) {
	resp, err := http.Get(fmt.Sprintf("http://api.fixer.io/%s?base=%s", when, base))
	if err != nil {
		return 0.00, fmt.Sprintf("No Response from Fixer: %v", err), http.StatusNoContent
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.00, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError
	}

	var fetchRates Rates
	err = json.Unmarshal(respBody, &fetchRates)
	if err != nil {
		return 0.00, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError
	}

	for k, v := range fetchRates.Rates {
		if k == target {
			return v, "", http.StatusOK
		}
	}

	return 0.00, "Not Found", http.StatusNotFound
}

// Rates : Structure to fetch latest rates into
type Rates struct {
	Base  string
	Rates map[string]float32
}
