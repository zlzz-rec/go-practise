package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewPrometheus(ctx context.Context) error {
	server := http.Server{Addr: "0.0.0.0:8399"}
	go func() {
		select {
		case <-ctx.Done():
			_ = server.Shutdown(ctx)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	return server.ListenAndServe()
}