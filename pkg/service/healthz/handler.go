// Package healthz contains handlers for readiness and liveness checks.
package healthz

import "net/http"

// NewHandler returns a new healthcheck handler.
func NewHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/readiness", healthReadinessHandler)
	mux.HandleFunc("/healthz/liveness", healthLivenessHandler)

	return mux
}
