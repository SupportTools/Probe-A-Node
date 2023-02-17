package run

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/supporttools/Probe-A-Node/modules/checks"
	"github.com/supporttools/Probe-A-Node/modules/cli"
	"github.com/supporttools/Probe-A-Node/modules/kubernetes"
	"github.com/supporttools/Probe-A-Node/modules/logging"
)

var log = logging.SetupLogging()

var (
	nodeStatus = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "node_status",
		Help: "The status of the node",
	})

	internalDns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "internal_dns",
		Help: "The status of the internal dns lookups",
	})

	externalDns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "external_dns",
		Help: "The status of the external dns lookups",
	})
)

type NodeResults struct {
	Timestamp   time.Time
	Status      bool
	KubeApi     bool
	InternalDns bool
	ExternalDns bool
}

// Run is the main function that is called when the program is elected leader
func Run() {
	nodeName := cli.Settings().NodeName
	log.Debugln("Node name: ", nodeName)
	client := kubernetes.SetupKubernetes()

	for {
		nodeResult := NodeResults{
			Timestamp: time.Now(),
		}

		log.Infoln("Getting node status")
		nodeInfo, err := kubernetes.GetNodeInfo(client, nodeName)
		if err != nil {
			nodeResult.Status = false
		}
		nodeResult.Status = nodeInfo.Status
		log.Infoln("Node status: ", nodeResult.Status)

		log.Infoln("Getting kube api info")
		version, err := checks.CheckKubeApi(client, nodeName)
		log.Info("Kube api version: ", version)
		if err != nil {
			nodeResult.KubeApi = false
		}
		nodeResult.KubeApi = true

		internalDnsEndpoint := cli.Settings().InternalDnsEndpoint
		log.Infoln("Checking internal dns info")
		if checks.CheckDns(internalDnsEndpoint, cli.Settings().InternalDnsServer) != nil {
			nodeResult.InternalDns = false
		}
		nodeResult.InternalDns = true

		externalDnsEndpoint := cli.Settings().ExternalDnsEndpoint
		log.Infoln("Checking external dns info")
		if checks.CheckDns(externalDnsEndpoint, cli.Settings().ExternalDnsServer) != nil {
			nodeResult.ExternalDns = false
		}
		nodeResult.ExternalDns = true

		log.Infoln("Checking overlay network")
		msg, unhealthyPods, err := checks.CheckOverlayNetwork(client, nodeName)
		if err != nil {
			log.Errorln("Error checking overlay network: ", err)
			log.Errorln("Message: ", msg)
		}
		log.Infoln("Unhealthy pods: ", unhealthyPods)
		//Need to add a check if the node is unhealthy and then set the node status to false

		log.Infoln("Sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}
}
