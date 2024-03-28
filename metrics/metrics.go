package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EVMVotesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "evm_votes_total",
			Help: "Number of EVM votes",
		},
		[]string{"network_name", "status"},
	)

	MaintainersGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "maintainers_status_list",
			Help: "Maintainer's status, 1 is active, 0 is not active",
		},
		[]string{"network_name"},
	)

	HeartbeatsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "heartbeats_total",
			Help: "Heartbeat of the application",
		},
		[]string{"status"},
	)
)
