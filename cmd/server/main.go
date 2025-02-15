package main

import (
	"github/erickmaria/go-api-observability/internal/config"
	"github/erickmaria/go-api-observability/internal/logger"
	"github/erickmaria/go-api-observability/internal/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	logger.NewLogger()
	config.NewConfig()
}

func main() {
	var port = config.GetSring("server.port")

	// create server mux
	mux := http.NewServeMux()
	registry := metrics.NewPrometheusRegistry()

	// create routes
	mux.HandleFunc("GET /", home)
	mux.Handle("GET /metrics", promhttp.Handler())

	// add prometheus middlewere
	prom := registry.Middleware(mux)

	// start server
	log.Info("server running on port ", port)
	if err := http.ListenAndServe(":"+port, prom); err != nil {
		log.Error("server error: ", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))
}
