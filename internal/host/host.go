package host

import (
	"context"
	"log/slog"
	"net/http"
)

const hostName = "127.0.0.1"

type WebHost struct {
	server http.Server
}

func New(port string, handler http.Handler) *WebHost {
	return &WebHost{
		server: http.Server{
			Addr:    hostName + ":" + port,
			Handler: handler,
		},
	}
}

func (h *WebHost) Run() {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			const errMsg = "error running web host"
			slog.With("error", err.Error()).Error(errMsg)
			panic(errMsg)
		}
	}()
}

func (h *WebHost) Terminate(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
