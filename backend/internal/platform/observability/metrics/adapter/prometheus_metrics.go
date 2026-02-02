package adapter

import (
	"dukkanim-api/internal/platform/observability/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusMetrics struct {
	registery *prometheus.Registry
}

func PrometheusMetrics() metrics.Metrics {
	reg := prometheus.NewRegistry()

	return &prometheusMetrics{
		registery: reg,
	}

}

func (p prometheusMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(p.registery, promhttp.HandlerOpts{})
}

func (p *prometheusMetrics) NewCounter(name, help string, labels []string) metrics.Counter {
	counterVec := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	p.registery.MustRegister(counterVec)

	return &prometheusCounter{
		counterVec: *counterVec,
	}
}

func (p *prometheusMetrics) NewHistogram(name, help string, labels []string, buckets []float64) metrics.Histogram {
	f_buckets := make([]float64, len(buckets))
	for i := range buckets {
		f_buckets[i] = float64(buckets[i])
	}

	histVec := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: f_buckets,
		},
		labels,
	)

	p.registery.MustRegister(histVec)

	return &prometheusHistogram{
		histVec: *histVec,
	}

}

func (p *prometheusMetrics) NewGauge(name, help string, labels []string) metrics.Gauge {

	// TODO: Cardinality improvements without labels!

	gaugeVec := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	p.registery.MustRegister(gaugeVec)

	return &prometheusGauge{
		gaugeVec: *gaugeVec,
	}
}

type prometheusCounter struct {
	counterVec prometheus.CounterVec
	counter    prometheus.Counter
}

func (p *prometheusCounter) Inc() {
	p.counter.Inc()
}
func (p *prometheusCounter) Add(add int) {
	p.counter.Add(float64(add))
}
func (p *prometheusCounter) WithLabelValues(labels map[string]string) metrics.Counter {
	return &prometheusCounter{
		counterVec: p.counterVec,
		counter:    p.counterVec.With(prometheus.Labels(labels)),
	}
}

type prometheusHistogram struct {
	histVec prometheus.HistogramVec
	hist    prometheus.Observer
}

func (p prometheusHistogram) Observe(value int) {
	p.hist.Observe(float64(value))
}

func (p *prometheusHistogram) WithLabelValues(labels map[string]string) metrics.Histogram {
	hist := p.histVec.With(prometheus.Labels(labels))
	p.hist = hist

	return p
}

type prometheusGauge struct {
	gaugeVec prometheus.GaugeVec
	gauge    prometheus.Gauge
}

func (p *prometheusGauge) Set(val int) {
	p.gauge.Set(float64(val))
}

func (p *prometheusGauge) Inc() {
	p.gauge.Inc()
}

func (p *prometheusGauge) Dec() {
	p.gauge.Dec()
}

func (p *prometheusGauge) Add(val int) {
	p.gauge.Add(float64(val))
}
func (p *prometheusGauge) Sub(val int) {
	p.gauge.Sub(float64(val))
}
func (p *prometheusGauge) SetToCurrentTime() {
	p.gauge.SetToCurrentTime()
}

func (p *prometheusGauge) WithLabelValues(labels map[string]string) metrics.Gauge {
	gauge := p.gaugeVec.With(prometheus.Labels(labels))
	p.gauge = gauge
	return p
}

// type linearBuckets struct {
// 	start float64
// 	width float64
// 	count int
// }

// func LinearBuckets(start, width float64, count int) metrics.HistogramBuckets {
// 	return linearBuckets{
// 		start: start,
// 		width: width,
// 		count: count,
// 	}
// }

// func (b linearBuckets) Define(start, width float64, count int) []float64 {
// 	return prometheus.LinearBuckets(b.start, b.width, b.count)
// }
