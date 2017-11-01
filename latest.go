package assignment2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/appengine"
	"github.com/golang/appengine/urlfetch"
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("http://api.fixer.io/latest")
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

	if payload.BaseCurrency == "EUR" || payload.TargetCurrency == "EUR" {
		for k, v := range latestRates.Rates {
			if k == payload.TargetCurrency {
				fmt.Fprint(w, v)
				return
			} else if k == payload.BaseCurrency {
				fmt.Fprint(w, 1/v)
				return
			}
		}
	} else {
		http.Error(w, "Only Euro conversion supported", http.StatusNotImplemented)
		return
	}

}

// Rates : Structre to fetc latest rates into
type Rates struct {
	Base  string
	Rates map[string]float32
}
