{{- .Workspace.TplHeader}}
// vim: set ft=go:

package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RequestMetrics struct {
	requestSize     *prometheus.SummaryVec
	responseSize    *prometheus.SummaryVec
	requestDuration *prometheus.HistogramVec
	requestCount    *prometheus.CounterVec
}

type RequestInfo struct {
	ServiceName string
	Hostname string
	URLPath string
}

func NewRequestMetrics() *RequestMetrics {
	// TODO: check error + fast check if already registered

	reqSize := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "service",
		Subsystem:  "api",
		Name:       "request_size_bytes",
		Help:       "Size of input request",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service", "host_name", "path"})
	prometheus.Register(reqSize)

	respSize := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "service",
		Subsystem:  "api",
		Name:       "response_size_bytes",
		Help:       "Size of response",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service", "host_name", "path"})
	prometheus.Register(respSize)

	reqDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "service",
		Subsystem: "api",
		Name:      "request_duration_seconds",
		Help:      "Duration of request in ms",
		Buckets:   []float64{0.02, 0.05, 0.1, 0.2, 0.5, 1, 5, 30, 60},
	}, []string{"service", "host_name", "path"})
	prometheus.Register(reqDuration)

	reqCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "service",
		Subsystem: "api",
		Name:      "requests_total",
		Help:      "Total count of processed requests",
	}, []string{"service", "host_name", "path", "status"})
	prometheus.Register(reqCount)

	return &RequestMetrics{
		requestSize:     reqSize,
		responseSize:    respSize,
		requestDuration: reqDuration,
		requestCount:    reqCount,
	}
}

func (rm *RequestMetrics) ReportRequestEnd(
	reqInfo RequestInfo,
	status int,
	duration time.Duration,
	requestSizeBytes int,
	responseSizeBytes int) {

	rm.requestSize.
		WithLabelValues(
			reqInfo.ServiceName,
			reqInfo.Hostname,
			reqInfo.URLPath).
		Observe(float64(requestSizeBytes))

	rm.responseSize.
		WithLabelValues(
			reqInfo.ServiceName,
			reqInfo.Hostname,
			reqInfo.URLPath).
		Observe(float64(responseSizeBytes))

	rm.requestDuration.
		WithLabelValues(
			reqInfo.ServiceName,
			reqInfo.Hostname,
			reqInfo.URLPath).
		Observe(duration.Seconds())

	rm.requestCount.
		WithLabelValues(
			reqInfo.ServiceName,
			reqInfo.Hostname,
			reqInfo.URLPath,
			strconv.Itoa(status)).
		Inc()
}
