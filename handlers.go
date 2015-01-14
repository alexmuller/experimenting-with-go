package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CspReport struct {
	Details CspDetails `json:"csp-report"`
}

type CspDetails struct {
	DocumentUri       string `json:"document-uri"`
	Referrer          string `json:"referrer"`
	BlockedUri        string `json:"blocked-uri"`
	ViolatedDirective string `json:"violated-directive"`
	OriginalPolicy    string `json:"original-policy"`
}

// JsonReceiverHandler receives JSON from a request body
// TODO: Validate and store the JSON
func JsonReceiverHandler(w http.ResponseWriter, req *http.Request) {
	var newCspReport CspReport

	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &newCspReport)

	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	w.Write([]byte("JSON received"))
}
