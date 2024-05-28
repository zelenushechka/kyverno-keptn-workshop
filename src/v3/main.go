package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := mux.NewRouter()

	// Register Prometheus metrics
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		},
		[]string{"path", "method"},
	)
	prometheus.MustRegister(requestDuration)

	// Register Request count metric
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Number of HTTP requests",
		},
		[]string{"path", "method"},
	)
	prometheus.MustRegister(requestCount)

	// Middleware to measure request duration
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start).Seconds()
			requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		})
	})

	// Middleware to measure request count
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
			next.ServeHTTP(w, r)
		})
	})

	// Expose health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Expose Prometheus metrics endpoint
	r.Path("/metrics").Handler(promhttp.Handler())

	// Route handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a new client
		client := openfeature.NewClient("app")
		featureEnabled, _ := client.BooleanValue(
			context.Background(), "my-new-feature", false, openfeature.EvaluationContext{},
		)

		// Delay the response between 5 and 10 seconds
		if featureEnabled {
			// return error 500 response
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Write([]byte("Hello World!"))
	})

	r.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		env := os.Environ()
		envMap := make(map[string]string)
		for _, e := range env {
			pair := strings.Split(e, "=")
			envMap[pair[0]] = pair[1]
		}
		json.NewEncoder(w).Encode(envMap)
	})

	// Initialize the FlagD client
	openfeature.SetProvider(flagd.NewProvider(
		flagd.WithHost("localhost"),
		flagd.FromEnv(),
	))

	// Start the server
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", r)
}
