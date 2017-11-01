package assignment2

import "net/http"

func avgHandler(w http.ResponseWriter, r *http.Request)  {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}