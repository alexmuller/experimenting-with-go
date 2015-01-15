package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"labix.org/v2/mgo"
)

var (
	mgoSession      *mgo.Session
	mgoDatabaseName = getenvDefault("JSON_REPORT_CATCHER_MONGO_DB", "json_report_catcher")
	mgoURL          = getenvDefault("JSON_REPORT_CATCHER_MONGO_URL", "localhost")
)

type CspReport struct {
	Details    CspDetails `json:"csp-report" bson:"csp_report"`
	ReportTime time.Time  `bson:"date_time"`
}

type CspDetails struct {
	DocumentUri       string `json:"document-uri" bson:"document_uri"`
	Referrer          string `json:"referrer" bson:"referrer"`
	BlockedUri        string `json:"blocked-uri" bson:"blocked_uri"`
	ViolatedDirective string `json:"violated-directive" bson:"violated_directive"`
	OriginalPolicy    string `json:"original-policy" bson:"original_policy"`
}

func getMgoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mgoURL)
		if err != nil {
			panic(err)
		}
	}
	return mgoSession.Clone()
}

func storeCspReport(report CspReport) {
	session := getMgoSession()
	defer session.Close()
	session.SetMode(mgo.Strong, true)

	collection := session.DB(mgoDatabaseName).C("reports")

	err := collection.Insert(report)

	if err != nil {
		panic(err)
	}

}

// JsonReceiverHandler receives JSON from a request body
// TODO: Validate the JSON
func JsonReceiverHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var newCspReport CspReport

	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &newCspReport)

	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	newCspReport.ReportTime = time.Now().UTC()

	go storeCspReport(newCspReport)

	w.Write([]byte("JSON received"))
}
