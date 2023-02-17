package health

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/supporttools/Probe-A-Node/modules/cli"
	"github.com/supporttools/Probe-A-Node/modules/logging"
)

var log = logging.SetupLogging()

var gitCommit string
var gitBranch string

func PrintVersion() {
	log.Printf("Current build version: %s", gitCommit)
	log.Printf("Current build branch: %s", gitBranch)
}

func StartMetricsServer() {
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/healthz", HealthHandler)
		router.HandleFunc("/version", VersionHandler)
		router.Handle("/metrics", promhttp.Handler())
		healthCheckPort := cli.Settings().HealthCheckPort
		if err := http.ListenAndServe(":"+healthCheckPort, router); err != nil {
			log.Fatal(err)
		} else {
			log.Infoln("Health server started")
		}
	}()
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Version: " + gitCommit))
}

func StartHealthServer() {
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/healthz", HealthHandler)
		router.HandleFunc("/version", VersionHandler)
		router.HandleFunc("/overall", OverallHandler)
		healthCheckPort := cli.Settings().HealthCheckPort
		if err := http.ListenAndServe(":"+healthCheckPort, router); err != nil {
			log.Fatal(err)
		} else {
			log.Infoln("Health server started")
		}
	}()
}

func OverallHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
