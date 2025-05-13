package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EVMVotesCounter = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "evm_votes_total",
			Help: "Number of EVM votes",
		},
		[]string{"chain", "address", "network", "status"},
	)

	MaintainersGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "maintainers_status_list",
			Help: "Maintainer's status, 1 is active, 0 is not active",
		},
		[]string{"chain", "address", "network"},
	)

	HeartbeatsCounter = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heartbeats_total",
			Help: "Heartbeat of the application",
		},
		[]string{"chain", "address", "status"},
	)

	chain   string
	address string
)

func Initialize(chain, address string) {
	chain = chain
	address = address
}

func SetEVMVotesMissed(targetNetwork string, missCnt int) {
	EVMVotesCounter.With(prometheus.Labels{"chain": chain, "address": address, "network": targetNetwork, "status": "missed"}).Add(float64(missCnt))
}

func SetEVMVotesSuccess(targetNetwork string, successCnt int) {
	EVMVotesCounter.With(prometheus.Labels{"chain": chain, "address": address, "network": targetNetwork, "status": "success"}).Add(float64(successCnt))
}

func SetMaintainersStatus(targetNetwork string, status int) {
	MaintainersGauge.With(prometheus.Labels{"chain": chain, "address": address, "network": targetNetwork}).Set(float64(status))
}

func SetHeartbeatsCounterMissed(missCnt int) {
	HeartbeatsCounter.With(prometheus.Labels{"chain": chain, "address": address, "status": "missed"}).Set(float64(missCnt))
}

func SetHeartbeatsCounterSuccess(successCnt int) {
	HeartbeatsCounter.With(prometheus.Labels{"chain": chain, "address": address, "status": "success"}).Set(float64(successCnt))
}
