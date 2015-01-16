package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alext/tablecloth"
)

var (
	reportingPort = getenvDefault("JSON_REPORT_CATCHER_REPORTING_PORT", ":8080")
)

func getenvDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}

	return val
}

func main() {
	if wd := os.Getenv("GOVUK_APP_ROOT"); wd != "" {
		tablecloth.WorkingDir = wd
	}

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/r", JsonReceiverHandler)

	log.Println("json-report-catcher: listening for reports on " + reportingPort)
	log.Fatal(tablecloth.ListenAndServe(reportingPort, publicMux, "reports"))
}
