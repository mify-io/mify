{{- .Workspace.TplHeader}}

package core

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ClientMetrics struct {
	requestDuration *prometheus.HistogramVec
	requestCount    *prometheus.CounterVec
}

func NewClientMetrics() *ClientMetrics {
	// TODO: check error + fast check if already registered

	reqDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "service",
		Subsystem: "api",
		Name:      "request_duration_seconds",
		Help:      "Duration of request in ms",
		Buckets:   []float64{0.02, 0.05, 0.1, 0.2, 0.5, 1, 5, 30, 60},
	}, []string{"service", "host_name", "target_service", "path"})
	prometheus.Register(reqDuration)

	reqCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "service",
		Subsystem: "api",
		Name:      "requests_total",
		Help:      "Total count of processed requests",
	}, []string{"service", "host_name", "target_service", "path", "status"})
	prometheus.Register(reqCount)

	return &ClientMetrics{
		requestDuration: reqDuration,
		requestCount:    reqCount,
	}
}

func (rm *ClientMetrics) ReportRequestEnd(
	reqCtx *MifyRequestContextBuilder,
	status int,
	duration time.Duration,
	targetService string,
	targetPath string) {

	rm.requestDuration.
		WithLabelValues(
			reqCtx.ServiceContext().ServiceName(),
			reqCtx.ServiceContext().Hostname(),
			targetService,
			targetPath).
		Observe(duration.Seconds())

	rm.requestCount.
		WithLabelValues(
			reqCtx.ServiceContext().ServiceName(),
			reqCtx.ServiceContext().Hostname(),
			targetService,
			targetPath,
			strconv.Itoa(status)).
		Inc()
}
