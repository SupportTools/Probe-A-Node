package main

import (
	"github.com/supporttools/Probe-A-Node/modules/cli"
	"github.com/supporttools/Probe-A-Node/modules/health"
	"github.com/supporttools/Probe-A-Node/modules/logging"
	"github.com/supporttools/Probe-A-Node/modules/run"
)

var log = logging.SetupLogging()

func main() {
	log.Infoln("Starting Probe-A-Node")
	health.PrintVersion()
	cli.Settings()
	logging.SetupLogging()
	log.Infoln("Starting Metrics Endpoint")
	health.StartMetricsServer()
	log.Infoln("Starting Health Endpoint")
	health.StartHealthServer()
	log.Infoln("Starting main process")
	run.Run()
}
