package socket

import (
	"github.com/prometheus/client_golang/prometheus"
)

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "socket",
		Subsystem:  "caro",
		Name:       "latency",
		Help:       "latency",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"api"},
)

func init() {
	prometheus.MustRegister(
		latency,
	)
}
