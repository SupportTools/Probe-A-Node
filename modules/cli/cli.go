package cli

import (
	"os"

	"github.com/supporttools/Probe-A-Node/modules/logging"
)

type Cli struct {
	HealthCheckPort     string
	Namespace           string
	PodName             string
	NodeName            string
	NodeInternalIP      string
	NodePublicIP        string
	InternalDnsEndpoint string
	ExternalDnsEndpoint string
	InternalDnsServer   string
	ExternalDnsServer   string
}

var log = logging.SetupLogging()

func Settings() Cli {
	healthCheckPort := os.Getenv("PORT")
	if healthCheckPort == "" {
		//log.Info("PORT is not set, defaulting to 9876")
		healthCheckPort = "9876"
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		log.Fatal("NAMESPACE is not set")
	}
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		log.Fatal("POD_NAME is not set")
	}
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		log.Fatal("NODE_NAME is not set")
	}
	nodeInternalIP := os.Getenv("NODE_INTERNAL_IP")
	if nodeInternalIP == "" {
		log.Fatal("NODE_INTERNAL_IP is not set")
	}
	nodePublicIP := os.Getenv("NODE_EXTERNAL_IP")
	if nodePublicIP == "" {
		log.Fatal("NODE_EXTERNAL_IP is not set")
	}
	internalDnsEndpoint := os.Getenv("INTERNAL_DNS_ENDPOINT")
	if internalDnsEndpoint == "" {
		log.Fatal("INTERNAL_DNS_ENDPOINT is not set")
	}
	externalDnsEndpoint := os.Getenv("EXTERNAL_DNS_ENDPOINT")
	if externalDnsEndpoint == "" {
		log.Fatal("EXTERNAL_DNS_ENDPOINT is not set")
	}
	internalDnsServer := os.Getenv("INTERNAL_DNS_SERVER")
	if internalDnsServer == "" {
		log.Fatal("INTERNAL_DNS_SERVER is not set")
	}
	externalDnsServer := os.Getenv("EXTERNAL_DNS_SERVER")
	if externalDnsServer == "" {
		log.Fatal("EXTERNAL_DNS_SERVER is not set")
	}

	settings := Cli{
		HealthCheckPort:     healthCheckPort,
		Namespace:           namespace,
		PodName:             podName,
		NodeName:            nodeName,
		NodeInternalIP:      nodeInternalIP,
		NodePublicIP:        nodePublicIP,
		InternalDnsEndpoint: internalDnsEndpoint,
		ExternalDnsEndpoint: externalDnsEndpoint,
		InternalDnsServer:   internalDnsServer,
		ExternalDnsServer:   externalDnsServer,
	}

	log.Debug(settings)
	return settings
}
