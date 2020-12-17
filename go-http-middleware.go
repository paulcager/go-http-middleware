package middleware

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeCachingHandler(age time.Duration, h http.Handler) http.Handler {
	ageSeconds := int64(math.Round(age.Seconds()))
	if ageSeconds <= 0 {
		return h
	}

	header := fmt.Sprintf("public,max-age=%d", ageSeconds)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", header)
			h.ServeHTTP(w, r)
		})
}

func MakeNoIndexHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Robots-Tag", "noindex")
			h.ServeHTTP(w, r)
		})
}

func MakeLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()

			uri := r.URL.String()
			method := r.Method
			fmt.Printf("%s %s %s %d\n", method, uri, r.RemoteAddr, end.Sub(start).Milliseconds())
		})
}

func EnablePrometheus() {
	http.Handle("/metrics", promhttp.Handler())
}

func EnablePrometheusForMux(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
