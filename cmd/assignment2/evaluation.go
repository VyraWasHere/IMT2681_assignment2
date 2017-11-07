package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

func evalHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := mgo.Dial(cfg.DBuri)
	defer sess.Close()
	if err != nil {
		http.Error(w, "No DB Connection", http.StatusInternalServerError)
		return
	}

	col := sess.DB(cfg.DBName).C("Webhooks")
	var res []Webhook

	err = col.Find(bson.M{}).All(&res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	notOkCount := 0
	for _, wh := range res {
		rate, _, status := getRate("latest", wh.BaseCurrency, wh.TargetCurrency)
		if status != http.StatusOK {
			notOkCount++
			continue
		}

		byteStr, err := wh.marshalInvoke(rate)
		if err != nil {
			notOkCount++
			continue
		}

		resp, err := http.Post(wh.WebhookURL, "application/json", bytes.NewBuffer(byteStr))
		if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent) {
			notOkCount++
			continue
		}
	}

	if notOkCount == 0 {
		fmt.Fprint(w, "Success")
	} else {
		fmt.Fprintf(w, "%v invokes failed", notOkCount)
	}
}

// Snippet from https://stackoverflow.com/a/31374980
func (wh *Webhook) marshalInvoke(currentRate float32) ([]byte, error) {
	return json.Marshal(struct {
		BaseCurrency   string  `json:"baseCurrency"`
		TargetCurrency string  `json:"targetCurrency"`
		CurrentRate    float32 `json:"currentRate"`
		MinTrigger     float32 `json:"minTriggerValue"`
		MaxTrigger     float32 `json:"maxTriggerValue"`
	}{wh.BaseCurrency, wh.TargetCurrency, currentRate, wh.MinTrigger, wh.MaxTrigger})
}
