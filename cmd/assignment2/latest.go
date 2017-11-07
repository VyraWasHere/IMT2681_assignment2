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

	resp, err := http.Get("http://api.fixer.io/latest?symbols=" + payload.BaseCurrency + "," + payload.TargetCurrency)
	if err != nil {
		http.Error(w, fmt.Sprintf("No Response from Fixer: %v", err), http.StatusNoContent)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var latestRates Rates
	err = json.Unmarshal(respBody, &latestRates)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for k, v := range latestRates.Rates {
		if k == payload.TargetCurrency {
			fmt.Fprint(w, v)
			return
		}
	}

}

// Rates : Structure to fetch latest rates into
type Rates struct {
	Base  string
	Rates map[string]float32
}
