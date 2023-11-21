{{- .TplHeader}}
// vim: set ft=go:

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MifyMetricsWrapper struct {
	register prometheus.Registerer
}

func (mfw *MifyMetricsWrapper) Register(collector prometheus.Collector) error {
	return mfw.register.Register(collector)
}

func NewMifyMetricsWrapper() *MifyMetricsWrapper {
	return &MifyMetricsWrapper{
		register: prometheus.DefaultRegisterer,
	}
}
