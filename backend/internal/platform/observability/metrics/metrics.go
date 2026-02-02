package metrics

import "net/http"

type Metrics interface {
	Handler() http.Handler
	NewCounter(name, help string, labels []string) Counter
	NewHistogram(name, help string, labels []string, buckets []float64) Histogram
	NewGauge(name, help string, labels []string) Gauge
}

type Counter interface {
	Inc()
	Add(int)
	WithLabelValues(map[string]string) Counter
}

type Gauge interface {
	Set(int)
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc()
	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
	// values.
	Dec()
	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(int)
	// Sub subtracts the given value from the Gauge. (The value can be
	// negative, resulting in an increase of the Gauge.)
	Sub(int)

	// SetToCurrentTime sets the Gauge to the current Unix time in seconds.
	SetToCurrentTime()

	WithLabelValues(map[string]string) Gauge
}

type Histogram interface {
	Observe(int)
	WithLabelValues(map[string]string) Histogram
}

// type HistogramBuckets interface {
// 	Define(start, width float64, count int) []float64
// }

func StandartLinearHistogram() []float64 {
	return []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0}
}
