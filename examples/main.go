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
