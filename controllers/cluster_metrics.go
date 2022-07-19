package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	metricNamespace = "upgrade_schedule_operator"
	metricSubsystem = "cluster"
)

// Counters for total applied and failed scheduled upgrades
var (
	counterLabels = []string{"cluster_id", "cluster_namespace", "origin_version", "target_version"}

	UpgradesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Subsystem: metricSubsystem,
			Name:      "scheduled_upgrades_applied_total",
			Help:      "Number of all scheduled upgrades applied",
		},
		counterLabels,
	)
	FailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Subsystem: metricSubsystem,
			Name:      "scheduled_upgrades_failed_total",
			Help:      "Number of all scheduled upgrades that failed to apply",
		},
		counterLabels,
	)
	SuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Subsystem: metricSubsystem,
			Name:      "scheduled_upgrades_succeeded_total",
			Help:      "Number of all scheduled upgrades that were applied successfully",
		},
		counterLabels,
	)
)

// Gauge for all clusters scheduled upgrade time
var (
	infoLabels = []string{"cluster_id", "cluster_namespace", "origin_version", "target_version"}

	UpgradesInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Subsystem: metricSubsystem,
			Name:      "scheduled_upgrades_time",
			Help:      "Gives upgrade time for all clusters.",
		},
		infoLabels,
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(UpgradesTotal, FailuresTotal, SuccessTotal, UpgradesInfo)
}
