package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	getRegExp, gErr  = regexp.Compile("^/exchange/[0-9a-fA-F]{16,32}/?$")
	postRegExp, pErr = regexp.Compile("^/exchange/?$")
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var (
		body   string
		status int
	)

	switch r.Method {
	case http.MethodPost:
		if pErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		match := postRegExp.Match([]byte(r.URL.Path))
		if !match {
			http.Error(w, "Please post new webhooks to /exchange", http.StatusBadRequest)
			return
		}

		body, status = processPost(r.Body)

	case http.MethodGet:
		if gErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		match := getRegExp.Match([]byte(r.URL.Path))
		if !match {
			http.Error(w, "Please provide a valid id", http.StatusBadRequest)
			return
		}

		body, status = processGet(strings.Split(r.URL.Path, "/")[2])

	case http.MethodDelete:
		if gErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		match := getRegExp.Match([]byte(r.URL.Path))
		if !match {
			http.Error(w, "Please provide a valid id", http.StatusBadRequest)
			return
		}

		body, status = processDelete(strings.Split(r.URL.Path, "/")[2])

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

func processPost(rBody io.ReadCloser) (body string, status int) {

	payload, err := ioutil.ReadAll(rBody)
	if err != nil {
		return http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError
	}

	var hook Webhook
	err = json.Unmarshal(payload, &hook)
	if err != nil {
		return http.StatusText(http.StatusBadRequest), http.StatusBadRequest
	}

	sess, err := mgo.Dial(cfg.DBuri)
	defer sess.Close()
	if err != nil {
		return "No DB Connection", http.StatusInternalServerError
	}

	col := sess.DB(cfg.DBName).C("Webhooks")
	info, err := col.Upsert(&hook, &hook)
	body = fmt.Sprintf("%v", info.UpsertedId)
	if info.UpsertedId == nil {
		status = http.StatusOK
	} else {
		status = http.StatusCreated
	}
	return body, status
}

func processGet(id string) (body string, status int) {
	sess, err := mgo.Dial(cfg.DBuri)
	defer sess.Close()
	if err != nil {
		return "No DB Connection", http.StatusInternalServerError
	}

	col := sess.DB(cfg.DBName).C("Webhooks")
	var res Webhook
	err = col.FindId(bson.ObjectIdHex(id)).One(&res)
	if err != nil {
		return fmt.Sprintf("Not found: %s", id), http.StatusNotFound
	}

	byteStr, err := json.Marshal(&res)
	if err != nil {
		return http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError
	}

	return string(byteStr), http.StatusOK
}

func newDBSession() (sess *mgo.Session, db *mgo.Database, err error) {
	if err == nil {
		db = sess.DB(cfg.DBName)
	}
	return
}

func processDelete(id string) (body string, status int) {
	sess, err := mgo.Dial(cfg.DBuri)
	defer sess.Close()
	if err != nil {
		return "No DB Connection", http.StatusInternalServerError
	}

	col := sess.DB(cfg.DBName).C("Webhooks")
	err = col.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return fmt.Sprintf("Not found: %s", id), http.StatusNotFound
	}

	return "", http.StatusOK
}

// Webhook : The payload structure
type Webhook struct {
	WebhookURL     string  `json:"webhookURL"`
	BaseCurrency   string  `json:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	MinTrigger     float32 `json:"minTriggerValue"`
	MaxTrigger     float32 `json:"maxTriggerValue"`
}
