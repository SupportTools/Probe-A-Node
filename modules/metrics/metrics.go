package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

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
	InternalDns bool
	ExternalDns bool
}
