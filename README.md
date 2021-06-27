# http-statsd

A Go's net/http middleware to emit metrics to StasD.

## Stats

It automatically creates metric name, then send the following stats
* counter (increment of 1)
* timing (response time in ms)
* sample rate of 1.0 (configurable)

### Metric Name

The middleware uses http.Request info for building up the metric name.

```golang
// where r is *http.Request
fmt.Sprintf("%s.%s.%d", r.URL.Path, r.Method, mt.Code)
```

Additionally, it converts url path like "/" into "." to make it friendly with StatsD backend service like Graphite. For example, `users/345/address` would be converted to `users.345.address`.

## Usage

The following example how to use this middleware with standard net/http.

```golang
package main

import (
	"net/http"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/gomwr/http-statsd/middleware"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

func main() {
	cfg := &statsd.ClientConfig{
		Address:     "127.0.0.1:8125",
		Prefix:      "http-statsd-example",
		UseBuffered: true,
	}
	cli, _ := statsd.NewClientWithConfig(cfg)

	mux := http.NewServeMux()
	mux.Handle("/", middleware.CaptureStats(cli)(http.HandlerFunc(hello)))

	http.ListenAndServe(":3000", mux)
}
```
