package main

import (
	"github/erickmaria/go-api-observability/internal/config"
	"github/erickmaria/go-api-observability/internal/metrics"
	routerand "github/erickmaria/go-api-observability/internal/routes/rand"
	"github/erickmaria/go-api-observability/internal/server"
	"log"
	"log/slog"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	config.NewConfig()
}

func main() {
	var port = config.GetSring("server.port")

	registry := metrics.NewPrometheusRegistry()

	// create server mux
	server := server.NewServer()
	// add middleweres
	server.Use(registry.Middleware)
	// create routes
	server.Handle("GET /rand", http.HandlerFunc(routerand.Random))
	server.Handle("GET /metrics", promhttp.Handler())
	// start server
	slog.Info("server running on port " + port)
	log.Fatal(server.ListenAndServe(":" + port))

}
