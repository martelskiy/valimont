package listener

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/martelskiy/valimont/internal/attestation"
	"github.com/martelskiy/valimont/internal/metric"
	"go.opentelemetry.io/otel"
)

type validator interface {
	GetAttestations(ctx context.Context) ([]attestation.Model, error)
	ValidatorsNr() uint16
}

type Listener struct {
	validator    validator
	pollInterval time.Duration
}

func New(validator validator, pollInterval time.Duration) *Listener {
	return &Listener{
		validator:    validator,
		pollInterval: pollInterval,
	}
}

func (l *Listener) Start(ctx context.Context) error {
	tracer := otel.Tracer("listener")
	ctx, span := tracer.Start(ctx, "Start")
	span.AddEvent("starting the listener")
	defer span.End()

	ticker := time.NewTicker(l.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			att, err := l.validator.GetAttestations(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					slog.Warn("received cancellation request. Cancelling listener...")
					return err
				}
				slog.Error("error received. Tolerating...")
				span.RecordError(err)
				continue
			}

			var attEpoch int
			if len(att) != 0 {
				attEpoch = att[0].Epoch
			}
			var attestationsSucceeded float64
			var attestationsMissed float64
			var inclusionDistanceTotal float64
			for _, a := range att {
				if a.Status == 1 {
					attestationsSucceeded++
					inclusionDistanceTotal += float64(a.InclusionDistance)
				} else {
					attestationsMissed++
				}
			}

			slog.
				With("attestationsNr", len(att)).
				With("attestationEpoch", attEpoch).
				With("attestationsMissed", attestationsMissed).
				With("attestationsSucceeded", attestationsSucceeded).
				With("avgInclusionDistance", inclusionDistanceTotal/attestationsSucceeded).
				Info("attestations fetched. Writing metric")

			metric.EpochTarget.Set(float64(attEpoch))
			metric.MissedAttestations.Set(float64(attestationsMissed))
			metric.SuccessAttestations.Set(float64(attestationsSucceeded))
			metric.AttestationRate.Set(float64(len(att)) / float64(l.validator.ValidatorsNr()))
			metric.InclusionDistance.Set(inclusionDistanceTotal / attestationsSucceeded)
		}
	}
}
