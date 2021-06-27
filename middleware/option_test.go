package middleware_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/gomwr/http-statsd/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {

	t.Run("create options with default values", func(t *testing.T) {
		emits := middleware.NewOptions()

		assert.Equal(t, float32(1.0), emits.SampleRate)
	})

	t.Run("create options with sampling rate 0.7", func(t *testing.T) {
		emits := middleware.NewOptionsWith(middleware.WithSample(0.7))

		assert.Equal(t, float32(0.7), emits.SampleRate)
	})

	t.Run("create options with sampling rate -1.0 then a warning log should display", func(t *testing.T) {
		out := CaptureOutput(func() {
			emits := middleware.NewOptionsWith(middleware.WithSample(-0.7))

			assert.Equal(t, float32(-0.7), emits.SampleRate)
		})

		assert.Contains(t, out, "warn: rate should be between 0.0 and 1.0")
		// Output: "2021/06/27 14:35:59 warn: rate should be between 0.0 and 1.0\n
	})
}

// CaptureOutput returns the standard log output as a string
// Taken from https://stackoverflow.com/a/26806093
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	f()

	log.SetOutput(os.Stderr)
	return buf.String()
}
