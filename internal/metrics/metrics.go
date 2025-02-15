package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests received",
	}, []string{"status", "path", "method"})

	activeRequestsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active connections to the service",
		},
	)

	latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"status", "path", "method"})

	// latencySummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	// 	Name: "http_request_duration_seconds",
	// 	Help: "Duration of requests summary in seconds",
	// 	Objectives: map[float64]float64{
	// 		0.5:  0.05,  // Median (50th percentile) with a 5% tolerance
	// 		0.9:  0.01,  // 90th percentile with a 1% tolerance
	// 		0.99: 0.001, // 99th percentile with a 0.1% tolerance
	// 	},
	// }, []string{"status", "path", "method"})
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

type registry struct{}

func NewPrometheusRegistry() *registry {
	prometheus.MustRegister(
		httpRequestCounter,
		activeRequestsGauge,
		latencyHistogram,
		// latencySummary,
	)

	return &registry{}
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (*registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		activeRequestsGauge.Inc()
		// Wrap the ResponseWriter to capture the status code
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(recorder, r)

		method := r.Method
		path := r.URL.Path // Path can be adjusted for aggregation (e.g., `/users/:id` → `/users/{id}`)
		status := strconv.Itoa(recorder.statusCode)
		// Increment the counter
		httpRequestCounter.WithLabelValues(status, path, method).Inc()
		//
		elapsed := time.Since(startTime).Seconds()
		latencyHistogram.WithLabelValues(status, path, method).Observe(elapsed)
		// latencySummary.WithLabelValues(status, path, method).Observe(elapsed)
		//
		activeRequestsGauge.Dec()
	})
}
