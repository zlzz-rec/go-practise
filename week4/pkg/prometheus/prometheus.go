package prometheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheus(port int, stop <-chan struct{}) error {
	server := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	go func() {
		<-stop
		server.Shutdown(context.Background())
	}()

	http.Handle("/metrics", promhttp.Handler())
	return server.ListenAndServe()
}
