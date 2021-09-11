package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	exporterCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "core",
			Subsystem: "caro",
			Name:      "api",
			Help:      "api counter",
		},
		[]string{"api"},
	)
	exporterLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "core",
			Subsystem:  "caro",
			Name:       "latency",
			Help:       "latency",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"api"},
	)
	exporterGame = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "core",
			Subsystem: "caro",
			Name:      "game",
			Help:      "game counter",
		},
		[]string{"label"},
	)
)

func init() {
	prometheus.MustRegister(
		exporterCounter,
		exporterLatency,
		exporterGame,
	)
}
