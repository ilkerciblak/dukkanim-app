package metrics

type HttpMetrics struct {
	Counter   Counter
	Histogram Histogram
	Gauge     Gauge
}

func HTTPMetrics(m Metrics) *HttpMetrics {

	return &HttpMetrics{
		Counter: m.NewCounter(
			"http_request_counter",
			"HTTP Request Hit Count",
			[]string{
				"http.method",
				"http.path",
				"http.response_status_code",
			},
		),
		Histogram: m.NewHistogram(
			"api_request_duration_ms",
			"Api Request Duration in Miliseconds",
			[]string{
				"http.method",
				"http.path",
			},
			StandartLinearHistogram(),
		),
		Gauge: m.NewGauge(
			"api_active_requests",
			"Active Request Count",
			[]string{
				"http_method",
				"http_path",
			},
		),
	}
}
