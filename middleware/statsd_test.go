package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gomwr/http-statsd/middleware"
	"github.com/gomwr/http-statsd/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCaptureStats(t *testing.T) {

	tests := []struct {
		name   string
		h      http.Handler
		url    string
		metric string
		rate   float32
	}{
		{
			"capture http stats",
			http.NotFoundHandler(),
			"http://api.example.com/countries/th/cities/cm",
			"countries.th.cities.cm.GET.404",
			1.0,
		},
		{
			"capture http stats with 0.7 sample rate",
			http.NotFoundHandler(),
			"http://api.example.com/countries/th/cities/cm",
			"countries.th.cities.cm.GET.404",
			0.7,
		},
		{
			"capture http stats with -0.2 sample rate",
			http.NotFoundHandler(),
			"http://api.example.com/countries/th/cities/cm",
			"countries.th.cities.cm.GET.404",
			-0.2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli := &mocks.Statter{}

			cli.On("Inc", mock.AnythingOfType("string"), int64(1), tt.rate).Return(nil)
			cli.On("Timing", mock.AnythingOfType("string"), mock.AnythingOfType("int64"), tt.rate).Return(nil)

			h := middleware.CaptureStats(cli)(http.NotFoundHandler())
			if tt.rate != 1.0 {
				h = middleware.CaptureStats(cli, middleware.WithSample(tt.rate))(http.NotFoundHandler())
			}
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)

			cli.AssertNumberOfCalls(t, "Inc", 1)
			cli.AssertCalled(t, "Inc", tt.metric, int64(1), tt.rate)
			cli.AssertNumberOfCalls(t, "Timing", 1)
			cli.AssertCalled(t, "Timing", tt.metric, mock.AnythingOfType("int64"), tt.rate)
		})
	}
}

func TestNoCaptureStats(t *testing.T) {
	t.Run("no stat capture", func(t *testing.T) {
		h := middleware.CaptureStats(nil)(http.NotFoundHandler())

		req := httptest.NewRequest(http.MethodGet, "https://api.example.com", nil)
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

func TestErrCaptureStats(t *testing.T) {
	t.Run("error sending data", func(t *testing.T) {
		cli := &mocks.Statter{}
		err := errors.New("failed sending data")
		rate := float32(1.0)
		url := "http://api.example.com/countries/th/cities/cm"
		metric := "countries.th.cities.cm.GET.404"

		out := CaptureOutput(func() {
			cli.On("Inc", mock.AnythingOfType("string"), int64(1), rate).Return(err)
			cli.On("Timing", mock.AnythingOfType("string"), mock.AnythingOfType("int64"), rate).Return(err)

			h := middleware.CaptureStats(cli)(http.NotFoundHandler())
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)

			cli.AssertNumberOfCalls(t, "Inc", 1)
			cli.AssertCalled(t, "Inc", metric, int64(1), rate)
			cli.AssertNumberOfCalls(t, "Timing", 1)
			cli.AssertCalled(t, "Timing", metric, mock.AnythingOfType("int64"), rate)
		})

		assert.Contains(t, out, err.Error())
	})
}
