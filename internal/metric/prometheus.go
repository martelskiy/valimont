package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EpochTarget = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "valimont",
		Name:      "epoch_target",
		Help:      "The validator epoch target number",
	})

	MissedAttestations = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "valimont",
		Name:      "missed_attestations",
		Help:      "How many validators missed",
	})

	SuccessAttestations = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "valimont",
		Name:      "success_attestations",
		Help:      "How many validations were successful",
	})

	AttestationRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "valimont",
		Name:      "attestation_rate",
		Help:      "How many validations were successful",
	})

	InclusionDistance = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "valimont",
		Name:      "inclusion_distance",
		Help:      "How many validations were successful",
	})
)
