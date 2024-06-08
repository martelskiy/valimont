package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/martelskiy/valimont/config"
	"github.com/martelskiy/valimont/internal/host"
	"github.com/martelskiy/valimont/internal/lifecycle"
	"github.com/martelskiy/valimont/internal/listener"
	"github.com/martelskiy/valimont/internal/metric"
	"github.com/martelskiy/valimont/internal/route"
	"github.com/martelskiy/valimont/internal/validator"
)

// @title valimont
// @version	1.0
func main() {
	ctx := context.Background()
	slog.Info("initializing otel collector")
	traceProvider, err := metric.InitializeOtel(ctx)
	if err != nil {
		panic(err)
	}

	router := route.NewRouter()
	router.
		WithAPIDocumentation().
		WithPrometheusMetrics()

	slog.With("port", config.Port).Info("runing web host")

	host := host.New(config.Port, router.GetRouter())
	host.Run()

	slog.Info("instantiating validator and listener")
	v := validator.New(config.ValidatorIndx, config.RateLimitPerMinute)
	listener := listener.New(v, config.PollInterval)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		slog.Info("starting to listen")
		if err := listener.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("all processes were terminated on cancellation context")
			} else {
				panic(err.Error())
			}
		}
	}()

	lifecycle.ListenForApplicationShutDown(func() {
		slog.Info("terminating the web host")
		defer func() {
			if err = traceProvider.Shutdown(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					slog.Info("trace provider was shut down on cancellaton context")
				} else {
					slog.With("err", err.Error()).Error("error shutting down tracer")
				}
			}
		}()
		defer cancel()
		if err := host.Terminate(ctx); err != nil {
			slog.With("err", err.Error()).Error("error during host termination")
		}

	}, make(chan os.Signal))
}
