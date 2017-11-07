package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	numOfDays = 7
)

func avgHandler(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()
	var sum float32
	sum = 0.0
	for i := 0; i < numOfDays; i++ {
		day := now.AddDate(0, 0, -i)
		rate, body, status := getRate(day.Format("2006-01-02"), payload.BaseCurrency, payload.TargetCurrency)
		if status != http.StatusOK {
			http.Error(w, body, status)
			return
		}
		sum += rate
	}

	fmt.Fprint(w, sum/numOfDays)
}
