package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/felixge/httpsnoop"
)

func CaptureStats(cli statsd.Statter, opts ...Option) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cli == nil {
				next.ServeHTTP(w, r)
				return
			}

			emits := NewOptionsWith(opts...)
			mt := httpsnoop.CaptureMetrics(next, w, r)
			// e.g. Metrics {Code:200 Duration:486.434274ms Written:56136}

			name := sanitizeName(fmt.Sprintf("%s.%s.%d", r.URL.Path, r.Method, mt.Code))
			if err := cli.Inc(name, 1, emits.SampleRate); err != nil {
				log.Printf("error sending stat: %s", err)
			}
			if err := cli.Timing(name, mt.Duration.Milliseconds(), emits.SampleRate); err != nil {
				log.Printf("error sending stat: %s", err)
			}
		})
	}
}

// sanitizeName tries formatting metric name to make it friendly with backend
// services such as Graphite
func sanitizeName(name string) string {
	name = strings.Trim(name, "/")
	name = strings.ReplaceAll(name, "/", ".")
	name = strings.ReplaceAll(name, ":", "")

	return name
}
