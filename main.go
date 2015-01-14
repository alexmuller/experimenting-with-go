package main

import (
	"log"
	"net/http"
	"os"
	"sync"

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

func catchListenAndServe(addr string, handler http.Handler, ident string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := tablecloth.ListenAndServe(addr, handler, ident)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if wd := os.Getenv("GOVUK_APP_ROOT"); wd != "" {
		tablecloth.WorkingDir = wd
	}

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/r", JsonReceiverHandler)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go catchListenAndServe(reportingPort, publicMux, "reports", wg)
	log.Println("json-report-catcher: listening for reports on " + reportingPort)

	wg.Wait()
}
