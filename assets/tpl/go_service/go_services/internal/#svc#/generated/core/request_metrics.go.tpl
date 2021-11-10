{{- .Workspace.TplHeader}}

package core

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RequestMetrics struct {
	requestDuration *prometheus.HistogramVec
}

func NewRequestMetrics() *RequestMetrics {
	reqDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "service",
		Subsystem: "api",
		Name:      "request_duration_seconds",
		Help:      "Duration of request in ms",
		Buckets:   []float64{0.02, 0.05, 0.1, 0.2, 0.5, 1, 5, 30, 60},
	}, []string{"host_name", "path"})
	prometheus.Register(reqDuration) // TODO: check error + fast check if already registered

	return &RequestMetrics{
		requestDuration: reqDuration,
	}
}

func (rm *RequestMetrics) ReportRequestEnd(reqCtx *MifyRequestContextBuilder, duration time.Duration) {
	rm.requestDuration.
		WithLabelValues(
			reqCtx.ServiceContext().Hostname(),
			reqCtx.GetURLPath()).
		Observe(duration.Seconds())
}
