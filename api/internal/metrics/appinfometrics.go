package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	portmetrics "github.com/ronnyp07/SportStream/api/internal/domain/ports/metrics"
)

type appInfoHandler struct {
	host      string
	version   string
	goVersion string
}

func NewAppInfoMetricsHandler(host, version, goVersion string) portmetrics.AppInfoMetricsHandler {
	return &appInfoHandler{
		host:      host,
		version:   version,
		goVersion: goVersion,
	}
}

func (m *appInfoHandler) RegisterVersionInfo() {
	versioninfo.WithLabelValues(m.host, m.version, m.goVersion).Inc()
}

var (
	versioninfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: MetricsNamespace,
		Name:      "version_info",
		Help:      "Version information.",
	}, []string{"host", "version", "go"})
)
